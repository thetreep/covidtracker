package graphql_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	ctx        context.Context
	httpClient *http.Client
	url        string
}

func NewClient(ctx context.Context, url string) *Client {
	return &Client{ctx: ctx, httpClient: http.DefaultClient, url: url}
}

func (c *Client) Do(template string, values map[string]interface{}) (*gqlResp, error) {

	//prepare request
	var payload bytes.Buffer
	body := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     template,
		Variables: values,
	}
	if err := json.NewEncoder(&payload).Encode(body); err != nil {
		return nil, errors.Wrap(err, "encode body")
	}

	req, err := http.NewRequest(http.MethodPost, c.url, &payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	for key, values := range req.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req = req.WithContext(c.ctx)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//format response
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		return nil, errors.Wrap(err, "reading body")
	}
	resp := &gqlResp{}
	if err := json.NewDecoder(&buf).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if len(resp.Errors) > 0 {
		// return first error
		return resp, resp.Errors[0]
	}
	return resp, nil
}

type gqlErr struct {
	Message string
}

func (e gqlErr) Error() string {
	return e.Message
}

type gqlResp struct {
	Data   map[string]interface{} `json:"data"`
	Errors []gqlErr
}
