package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	nhttp "net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/longbridgeapp/openapi-go/config"
	"github.com/longbridgeapp/openapi-go/log"
)

type apiResponse struct {
	Code    int
	Message string
	Data    json.RawMessage
}

type otpResponse struct {
	Otp string
}

// Client is a http client to access Longbridge REST OpenAPI
type Client struct {
	opts       *Options
	httpClient *nhttp.Client
}

// Get sends Get request with queryParams
func (c *Client) Get(ctx context.Context, path string, queryParams url.Values, resp interface{}) error {
	return c.Call(ctx, "GET", path, queryParams, nil, resp)
}

// Post sends Post request with json body
func (c *Client) Post(ctx context.Context, path string, body interface{}, resp interface{}) error {
	return c.Call(ctx, "POST", path, nil, body, resp)
}

// Put sends Put request with json body
func (c *Client) Put(ctx context.Context, path string, body interface{}, resp interface{}) error {
	return c.Call(ctx, "PUT", path, nil, body, resp)
}

// Delete sends Delete request with queryParams
func (c *Client) Delete(ctx context.Context, path string, queryParams interface{}, resp interface{}) error {
	return c.Call(ctx, "DELETE", path, queryParams, nil, resp)
}

// GetOTP to get one time password
// Reference: https://open.longbridgeapp.com/en/docs/socket-token-api
func (c *Client) GetOTP(ctx context.Context) (string, error) {
	res := &otpResponse{}
	err := c.Get(ctx, "/v1/socket/token", nil, res)
	if err != nil {
		return "", err
	}
	return res.Otp, nil
}

func (c *Client) GetOTPV2(ctx context.Context) (string, error) {
	res := &otpResponse{}
	err := c.Get(ctx, "/v2/socket/token", nil, res)
	if err != nil {
		return "", err
	}
	return res.Otp, nil
}

// Call will send request with signature to http server
func (c *Client) Call(ctx context.Context, method, path string, queryParams interface{}, body interface{}, resp interface{}) (err error) {
	var (
		br       io.Reader
		bb       []byte
		httpResp *nhttp.Response
		rb       []byte
	)

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

	if queryParams != nil {
		vals, ok := queryParams.(url.Values)
		if !ok {
			if vals, err = query.Values(queryParams); err != nil {
				return
			}
		}
		req.URL.RawQuery = vals.Encode()
	}
	req.Header.Add("x-api-key", c.opts.AppKey)
	req.Header.Add("authorization", c.opts.AccessToken)
	if len(bb) != 0 {
		req.Header.Add("content-type", "application/json; charset=utf-8")
	}
	signature(req, c.opts.AppSecret, bb)

	log.Debugf("http call method:%v url:%v body:%v", req.Method, req.URL, string(bb))
	req.Close = true
	httpResp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	log.Debugf("http call response headers:%v", httpResp.Header)

	defer httpResp.Body.Close()
	rb, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	log.Debugf("http call response body:%v", string(rb))
	apiResp := &apiResponse{}

	if isJSON(httpResp.Header.Get("content-type")) {
		if err = json.Unmarshal(rb, apiResp); err != nil {
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

	if err = json.Unmarshal(apiResp.Data, resp); err != nil {
		return err
	}
	return nil
}

func isJSON(ct string) bool {
	return strings.Contains(ct, "application/json")
}

// New create http client to call Longbridge REST OpenAPI
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
	return New(
		WithAccessToken(c.AccessToken),
		WithAppKey(c.AppKey),
		WithAppSecret(c.AppSecret),
		WithTimeout(c.HTTPTimeout),
		WithClient(c.Client),
		WithURL(c.HttpURL),
	)
}
