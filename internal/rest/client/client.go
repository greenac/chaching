package rest

import (
	"bytes"
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/rest/models"
	"io"
	"net/http"
)

var _ models.IClient = (*Client)(nil)

type Client struct {
	BaseHeaders *models.Headers
	HttpClient  models.IHttpClient
	BodyReader  func(r io.Reader) ([]byte, error)
	GetRequest  func(method, url string, body io.Reader) (*http.Request, error)
}

func (c *Client) Get(url string, headers *models.Headers, params models.UrlParams) (models.Response, genErr.IGenError) {
	req, err := c.GetRequest("GET", url, nil)
	if err != nil {
		ge := genErr.GenError{}
		return models.Response{}, ge.AddMsg("ClientImpl:Get:failed to get new request with url: " + url + "with error: " + err.Error())
	}

	return c.makeRequest(req, c.makeHeaders(headers))
}

func (c *Client) PostBody(url string, headers *models.Headers, body []byte) (models.Response, genErr.IGenError) {
	req, err := c.GetRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		ge := genErr.GenError{}
		return models.Response{}, ge.AddMsg("ClientImpl:PostBody:failed to get new request with url: " + url + "with error: " + err.Error())
	}

	return c.makeRequest(req, c.makeHeaders(headers))
}

func (c *Client) PostUrl(url string, headers *models.Headers, params models.UrlParams) (models.Response, genErr.IGenError) {
	req, err := c.GetRequest("POST", url, nil)
	if err != nil {
		ge := genErr.GenError{}
		return models.Response{}, ge.AddMsg("ClientImpl:PostUrl:failed to get new request with url: " + url + "with error: " + err.Error())
	}

	return c.makeRequest(req, c.makeHeaders(headers))
}

func (c *Client) makeHeaders(headers *models.Headers) map[string]string {
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

func (c *Client) makeRequest(req *http.Request, headers map[string]string) (models.Response, genErr.IGenError) {
	for k, h := range headers {
		req.Header.Add(k, h)
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		ge := genErr.GenError{}
		return models.Response{}, ge.AddMsg("ClientImpl:makeRequest:failed to make request with error: " + err.Error())
	}

	var body []byte
	if res.Body != nil {
		b, err := c.BodyReader(res.Body)
		if err != nil {
			ge := genErr.GenError{}
			return models.Response{}, ge.AddMsg("ClientImpl:makeRequest:failed to read response body with error: " + err.Error())
		}

		body = b
	}

	return models.Response{StatusCode: res.StatusCode, Status: res.Status, Body: body}, nil
}
