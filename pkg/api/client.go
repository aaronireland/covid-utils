package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type responseReader func(*http.Response, interface{}) error

type Client struct {
	BaseURL      string
	HTTPClient   *http.Client
	accept       string
	contentType  string
	readResponse responseReader
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (p Params) Encode() string {
	values := url.Values{}
	for _, value := range p {
		values.Set(value.Key, value.Value)
	}
	return values.Encode()
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	client := &Client{
		BaseURL:      baseURL,
		HTTPClient:   &http.Client{Timeout: time.Minute},
		readResponse: defaultReadResponse,
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			panic(fmt.Sprintf("API Client Option function failed: %v", err))
		}
	}

	return client
}

func (c *Client) URI(route string, params ...Param) string {
	var path string
	if len(route) < 2 {
		path = c.BaseURL
	} else if strings.HasPrefix(route, "/") {
		path = fmt.Sprintf("%s%s", c.BaseURL, route)
	} else {
		path = fmt.Sprintf("%s/%s", c.BaseURL, route)
	}

	if len(params) > 0 {
		filters := url.Values{}
		for _, param := range params {
			filters.Add(param.Key, param.Value)
		}
		path += ("?" + filters.Encode())
	}

	return path
}

func (c *Client) Do(req *http.Request, resp interface{}) error {
	if c.accept != "" {
		req.Header.Set("Accept", c.accept)
	}

	if c.contentType != "" {
		req.Header.Set("Content-Type", c.contentType)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to get response: %s", err)
	}

	return c.readResponse(res, resp)

}
