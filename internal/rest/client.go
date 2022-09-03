package rest

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	genErr "github.com/greenac/chaching/internal/error"
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
	Get(url string, headers *Headers, params UrlParams) (Response, *genErr.GenError)
	PostBody(url string, headers *Headers, body []byte) (Response, *genErr.GenError)
	PostUrl(url string, headers *Headers, params UrlParams) (Response, *genErr.GenError)
}

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var _ IClient = (*ClientImpl)(nil)

type ClientImpl struct {
	BaseHeaders *Headers
	HttpClient  IHttpClient
	BodyReader  func(r io.Reader) ([]byte, error)
	GetRequest  func(method, url string, body io.Reader) (*http.Request, error)
}

func (c *ClientImpl) Get(url string, headers *Headers, params UrlParams) (Response, *genErr.GenError) {
	req, err := c.GetRequest("GET", url, nil)
	if err != nil {
		ge := genErr.GenError{}
		return Response{}, ge.AddMsg("ClientImpl:Get:failed to get new request with url: " + url + "with error: " + err.Error())
	}

	return c.makeRequest(req, c.makeHeaders(headers))
}

func (c *ClientImpl) PostBody(url string, headers *Headers, body []byte) (Response, *genErr.GenError) {
	req, err := c.GetRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		ge := genErr.GenError{}
		return Response{}, ge.AddMsg("ClientImpl:PostBody:failed to get new request with url: " + url + "with error: " + err.Error())
	}

	return c.makeRequest(req, c.makeHeaders(headers))
}

func (c *ClientImpl) PostUrl(url string, headers *Headers, params UrlParams) (Response, *genErr.GenError) {
	req, err := c.GetRequest("POST", url, nil)
	if err != nil {
		ge := genErr.GenError{}
		return Response{}, ge.AddMsg("ClientImpl:PostUrl:failed to get new request with url: " + url + "with error: " + err.Error())
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

func (c *ClientImpl) makeRequest(req *http.Request, headers map[string]string) (Response, *genErr.GenError) {
	for k, h := range headers {
		req.Header.Add(k, h)
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		ge := genErr.GenError{}
		return Response{}, ge.AddMsg("ClientImpl:makeRequest:failed to make request with error: " + err.Error())
	}

	var body []byte
	if res.Body != nil {
		b, err := c.BodyReader(res.Body)
		if err != nil {
			ge := genErr.GenError{}
			return Response{}, ge.AddMsg("ClientImpl:makeRequest:failed to read response body with error: " + err.Error())
		}

		body = b
	}

	return Response{StatusCode: res.StatusCode, Status: res.Status, Body: body}, nil
}
