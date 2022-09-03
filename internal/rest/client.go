package rest

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type HeaderValue []string

func (hv HeaderValue) String() string {
	return strings.Join(hv, ", ")
}

type Headers map[string]HeaderValue

type UrlParamType string

const (
	UrlParamTypeEqual              UrlParamType = "="
	UrlParamTypeLessThan           UrlParamType = "<"
	UrlParamTypeLessThanEqualTo    UrlParamType = "<="
	UrlParamTypeGreaterThan        UrlParamType = ">"
	UrlParamTypeGreaterThanEqualTo UrlParamType = ">="
)

type UrlParam struct {
	LVal    string
	RVal    string
	Compare UrlParamType
}

func (up UrlParam) String() string {
	return up.LVal + string(up.Compare) + up.RVal
}

type UrlParams []UrlParam

func (ups UrlParams) String() string {
	s := ""
	for i, p := range ups {
		s += p.String()
		if i < len(ups)-1 {
			s += "?"
		}
	}

	return s
}

type Response struct {
	StatusCode int
	Status     string
	Body       []byte
}

type IClient interface {
	Get(url string, headers *Headers, params UrlParams) (Response, error)
	PostBody(url string, headers *Headers, body []byte) (Response, error)
	PostUrl(url string, headers *Headers, params UrlParams) (Response, error)
}

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var _ IClient = (*ClientImpl)(nil)

type ClientImpl struct {
	BaseHeaders *Headers
	HttpClient  IHttpClient
	BodyReader  func(r io.Reader) ([]byte, error)
}

func (c *ClientImpl) Get(url string, headers *Headers, params UrlParams) (Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, err
	}

	return c.makeRequest(req, c.makeHeaders(headers))
}

func (c *ClientImpl) PostBody(url string, headers *Headers, body []byte) (Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}

	return c.makeRequest(req, c.makeHeaders(headers))
}

func (c *ClientImpl) PostUrl(url string, headers *Headers, params UrlParams) (Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return Response{}, err
	}

	return c.makeRequest(req, c.makeHeaders(headers))
}

func (c *ClientImpl) makeHeaders(headers *Headers) map[string]string {
	reqHeaders := map[string]string{}

	if c.BaseHeaders != nil {
		for key, hds := range *c.BaseHeaders {
			reqHeaders[key] = hds.String()
		}
	}

	if headers != nil {
		for key, hds := range *headers {
			reqHeaders[key] = hds.String()
		}
	}

	return reqHeaders
}

func (c *ClientImpl) makeRequest(req *http.Request, headers map[string]string) (Response, error) {
	for k, h := range headers {
		req.Header.Add(k, h)
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return Response{}, err
	}

	var body []byte
	if res.Body != nil {
		b, err := c.BodyReader(res.Body)
		if err != nil {
			return Response{}, err
		}

		body = b
	}

	return Response{StatusCode: res.StatusCode, Status: res.Status, Body: body}, nil
}
