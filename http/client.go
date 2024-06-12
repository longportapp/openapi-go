package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	nhttp "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"

	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/log"
)

const (
	headerRetryAfter = "retry-after"
	headerTraceID    = "x-trace-id"
)

type apiResponse struct {
	Code    int
	Message string
	Data    json.RawMessage
	TraceID string
}

type otpResponse struct {
	Otp string
}

var _ error = (*retryError)(nil)

type retryError struct{}

func (re *retryError) Error() string {
	return "retry"
}

var errNeedRetry = &retryError{}

// Client is a http client to access Longbridge REST OpenAPI
type Client struct {
	opts       *Options
	httpClient *nhttp.Client
}

// RequestOptions use to set additional information for the request
type RequestOptions struct {
	// Request Header
	Header nhttp.Header
	body   interface{}
}

// RequestOption use to set addition info to request
type RequestOption func(*RequestOptions)

// WithHeader set request header
func WithHeader(h nhttp.Header) RequestOption {
	return func(o *RequestOptions) {
		o.Header = h
	}
}

// WithBody to set playload
func WithBody(v interface{}) RequestOption {
	return func(o *RequestOptions) {
		if v != nil {
			o.body = v
		}
	}
}

// Get sends Get request with queryParams
func (c *Client) Get(ctx context.Context, path string, queryParams url.Values, resp interface{}, ropts ...RequestOption) error {
	return c.Call(ctx, "GET", path, queryParams, nil, resp, ropts...)
}

// Post sends Post request with json body
func (c *Client) Post(ctx context.Context, path string, body interface{}, resp interface{}, ropts ...RequestOption) error {
	return c.Call(ctx, "POST", path, nil, body, resp, ropts...)
}

// Put sends Put request with json body
func (c *Client) Put(ctx context.Context, path string, body interface{}, resp interface{}, ropts ...RequestOption) error {
	return c.Call(ctx, "PUT", path, nil, body, resp, ropts...)
}

// Delete sends Delete request with queryParams
func (c *Client) Delete(ctx context.Context, path string, queryParams interface{}, resp interface{}, ropts ...RequestOption) error {
	return c.Call(ctx, "DELETE", path, queryParams, nil, resp, ropts...)
}

// GetOTP to get one time password
// Reference: https://open.longportapp.com/en/docs/socket-token-api
func (c *Client) GetOTP(ctx context.Context, ropts ...RequestOption) (string, error) {
	res := &otpResponse{}
	err := c.Get(ctx, "/v1/socket/token", nil, res, ropts...)
	if err != nil {
		return "", err
	}
	return res.Otp, nil
}

func (c *Client) GetOTPV2(ctx context.Context, ropts ...RequestOption) (string, error) {
	res := &otpResponse{}
	err := c.Get(ctx, "/v2/socket/token", nil, res, ropts...)
	if err != nil {
		return "", err
	}
	return res.Otp, nil
}

// Call will send request with signature to http server
func (c *Client) Call(ctx context.Context, method, path string, queryParams interface{}, body interface{}, resp interface{}, ropts ...RequestOption) (err error) {
	var (
		br       io.Reader
		bb       []byte
		httpResp *nhttp.Response
		rb       []byte
	)

	ro := &RequestOptions{}
	for _, opt := range ropts {
		opt(ro)
	}

	if body == nil && ro.body != nil {
		body = ro.body
	}

	if body != nil {
		bb, err = json.Marshal(body)
		if err != nil {
			return err
		}
		br = bytes.NewBuffer(bb)
	}

	req, err := nhttp.NewRequestWithContext(ctx, method, c.opts.URL+path, br)
	if err != nil {
		return err
	}

	if ro.Header != nil {
		for k, v := range ro.Header {
			req.Header[k] = v
		}
	}

	if queryParams != nil {
		vals, ok := queryParams.(url.Values)
		if !ok {
			if vals, err = query.Values(queryParams); err != nil {
				return
			}
		}
		req.URL.RawQuery = vals.Encode()
	}
	req.Header.Set("x-api-key", c.opts.AppKey)
	req.Header.Set("authorization", c.opts.AccessToken)
	if len(bb) != 0 {
		req.Header.Set("content-type", "application/json; charset=utf-8")
	}

	log.Debugf("http call method:%v url:%v body:%s", req.Method, req.URL, bb)
	var retryCount int
	call := func(disabeRetry bool) error {
		if retryCount > 0 {
			log.Debugf("retry calling %s %s, count: %d", req.Method, req.URL.Path, retryCount)
			// reset x-timestamp in signature func
			req.Header.Set(headerTimestamp, "")
		}

		signature(req, c.opts.AppSecret, bb)

		httpResp, err = c.httpClient.Do(req)
		if err != nil {
			return err
		}

		// if disabled retry just return
		if disabeRetry {
			return nil
		}

		wait, ok := parseRatelimit(httpResp)

		if !ok {
			// if not rate limited just return
			return nil
		}

		// need retry so close body
		httpResp.Body.Close()
		time.Sleep(wait)
		return errNeedRetry
	}

	for {
		err = call(c.opts.DisableRetry)

		if err == nil {
			break
		}

		if errors.Is(err, errNeedRetry) {
			retryCount = retryCount + 1
			// retry
			continue

		}
		// handle error
		break
	}

	log.Debugf("http call response headers:%v", httpResp.Header)
	defer httpResp.Body.Close()
	if rb, err = io.ReadAll(httpResp.Body); err != nil {
		return err
	}
	log.Debugf("http call response body:%s", rb)

	apiResp := &apiResponse{}

	if v := httpResp.Header.Get("x-trace-id"); v != "" {
		apiResp.TraceID = v
	}

	if isJSON(httpResp.Header.Get("content-type")) {
		if err = jsonUnmarshal(bytes.NewReader(rb), apiResp); err != nil {
			return err
		}
	} else {
		apiResp.Message = string(rb)
	}

	if httpResp.StatusCode != nhttp.StatusOK || apiResp.Code != 0 {
		return NewError(httpResp.StatusCode, apiResp)
	}

	if resp == nil {
		return
	}

	if err = jsonUnmarshal(bytes.NewReader(apiResp.Data), resp); err != nil {
		return err
	}
	return nil
}

func isJSON(ct string) bool {
	return strings.Contains(ct, "application/json")
}

func parseRatelimit(res *nhttp.Response) (time.Duration, bool) {
	if res.StatusCode != http.StatusTooManyRequests {
		return 0, false
	}

	waitStr := res.Header.Get(headerRetryAfter)
	if waitStr == "" {
		return 0, false
	}

	d, err := time.ParseDuration(waitStr)
	if err != nil {
		return 0, false
	}

	log.Warnf("request %s %s (%s) has been rate limited, will retry after %v...", res.Request.Method, res.Request.URL.Path, res.Header.Get(headerTraceID), d)

	return d, true
}

func jsonUnmarshal(r io.Reader, v interface{}) error {
	d := json.NewDecoder(r)
	d.UseNumber()
	return d.Decode(v)
}

// New create http client to call Longbridge REST OpenAPI
//
// Example:
//	cli, err := New(http.WithAppKey("appkey"), http.WithAppSecret("appSecret"), http.WithAccessToken("token"))
//  cli.Do

func New(opt ...Option) (*Client, error) {
	opts := newOptions(opt...)
	if opts.URL == "" {
		return nil, errors.New("http url is empty")
	}

	cli := &nhttp.Client{Timeout: opts.Timeout}

	if opts.Client != nil {
		cli = opts.Client
	}

	client := &Client{
		opts:       opts,
		httpClient: cli,
	}
	return client, nil
}

// NewFromCfg init longbridge http client from *config.Config
func NewFromCfg(c *config.Config) (*Client, error) {
	cli, err := New(
		WithAccessToken(c.AccessToken),
		WithAppKey(c.AppKey),
		WithAppSecret(c.AppSecret),
		WithTimeout(c.HTTPTimeout),
		WithClient(c.Client),
		WithURL(c.HttpURL),
	)
	if err != nil {
		return cli, err
	}

	if c.DisalbeHTTPRetry {
		cli.opts.DisableRetry = true
	}
	return cli, nil
}
