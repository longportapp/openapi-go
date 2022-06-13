package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	nhttp "net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)


type ApiResponse struct {
	code int
	message string
	data json.RawMessage
}

type Client struct {
	opts *Options
	httpClient *nhttp.Client
}

func (c *Client) Get(ctx context.Context, path string, queryParams interface{}, resp interface{}) error {
	return c.Call(ctx, "GET", path, queryParams, nil, resp)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, resp interface{}) error {
	return c.Call(ctx, "POST", path, nil, body, resp)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}, resp interface{}) error {
	return c.Call(ctx, "PUT", path, nil, body, resp)
}

func (c *Client) Delete(ctx context.Context, path string, queryParams interface{}, resp interface{}) error {
	return c.Call(ctx, "DELETE", path, queryParams, nil, resp)
}

func (c *Client) Call(ctx context.Context, method, path string, queryParams interface{}, body interface{}, resp interface{}) (err error) {
	var (
		br io.Reader
	    bb []byte
		httpResp *nhttp.Response
		rb []byte
	)

	if body != nil {
		bb, err = json.Marshal(body)
		if err != nil {
			return err
		}
		br = bytes.NewReader(bb)
	}
	req, err := nhttp.NewRequestWithContext(ctx, method, c.opts.URL + path, br)
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
		req.URL.RawQuery =  vals.Encode()
	}
	req.Header.Add("X-Api-Key", c.opts.AppKey)
	req.Header.Add("Authorization", c.opts.AccessToken)

	signature(req, c.opts.AppSecret, bb)

	httpResp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer httpResp.Body.Close()
	rb, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	apiResp := &ApiResponse{}
	if err = json.Unmarshal(rb, apiResp); err != nil {
		return err
	}

	if httpResp.StatusCode != nhttp.StatusOK || apiResp.code != 0 {
		return NewError(httpResp.StatusCode, apiResp)
	}

	if resp == nil {
		return
	}
	if err = json.Unmarshal(apiResp.data, resp); err != nil {
		return err
	}
	return nil
}

