package http

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var json = jsoniter.ConfigFastest

type (
	Client struct {
		host     string
		provider fasthttp.Client
	}

	Request struct {
		Path    string
		Method  string
		Headers map[string]string
		Body    interface{}
	}

	Response struct {
		StatusCode int
		Out        any
		OutError   error
	}
)

func (c *Client) Do(req *Request, res *Response) error {
	httpReq := fasthttp.AcquireRequest()
	httpReq.SetRequestURI(fmt.Sprintf("%s%s", c.host, req.Path))
	httpReq.Header.SetMethod(req.Method)

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return err
		}

		httpReq.SetBody(bodyBytes)
	}

	httpRes := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpRes)

	err := c.provider.Do(httpReq, httpRes)
	fasthttp.ReleaseRequest(httpReq)

	if err != nil {
		return ApiError{
			Message: err.Error(),
		}
	}

	if res == nil {
		return nil
	}

	if httpRes.StatusCode() != res.StatusCode {
		err = json.Unmarshal(httpRes.Body(), res.OutError)
		if err == nil {
			return res.OutError
		}
	} else if res.Out != nil {
		err = json.Unmarshal(httpRes.Body(), res.Out)
	}

	if err != nil {
		return ApiError{
			Message: err.Error(),
		}
	}

	return nil
}

func NewCient(host string) *Client {
	return &Client{
		host: host,
		provider: fasthttp.Client{
			ReadTimeout:                   time.Second * 1,
			MaxIdleConnDuration:           time.Second * 10,
			NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
			DisableHeaderNamesNormalizing: true,
			DisablePathNormalizing:        true,
			Dial: (&fasthttp.TCPDialer{
				Concurrency:      4096,
				DNSCacheDuration: time.Hour * 1,
			}).Dial,
		},
	}
}
