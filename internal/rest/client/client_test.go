package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/rest/mocks"
	"github.com/greenac/chaching/internal/rest/models"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type mockRespObj struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       string `json:"age"`
}

func TestHeaderValue_String(t *testing.T) {
	Convey("HeaderValue_String", t, func() {
		hv := models.HeaderValue{"one", "two", "three"}
		So(hv.String(), ShouldEqual, "one, two, three")
	})
}

func TestUrlParam_String(t *testing.T) {
	Convey("UrlParam_String", t, func() {
		param := models.UrlParam{LVal: "beach", Compare: models.UrlParamTypeEqual, RVal: "ball"}
		So(param.String(), ShouldEqual, "beach=ball")
	})
}

func TestUrlParams_String(t *testing.T) {
	Convey("UrlParams_String", t, func() {
		params := models.UrlParams{
			models.UrlParam{LVal: "beach", Compare: models.UrlParamTypeEqual, RVal: "ball"},
			models.UrlParam{LVal: "stevie", Compare: models.UrlParamTypeGreaterThanEqualTo, RVal: "wonder"},
			models.UrlParam{LVal: "style", Compare: models.UrlParamTypeLessThan, RVal: "3"},
		}
		So(params.String(), ShouldEqual, "beach=ball?stevie>=wonder?style<3")
	})
}

func TestClientImpl_makeHeaders(t *testing.T) {
	Convey("ClientImpl_makeHeaders", t, func() {
		Convey("ClientImpl_makeHeaders with base headers", func() {
			bh := models.Headers{
				"one": models.HeaderValue{"a"},
				"two": models.HeaderValue{"b, c"},
			}
			c := Client{BaseHeaders: &bh}
			headers := c.makeHeaders(nil)
			So(headers, ShouldResemble, map[string]string{"one": "a", "two": "b, c"})
		})

		Convey("ClientImpl_makeHeaders with headers param", func() {
			bh := models.Headers{
				"one": models.HeaderValue{"a"},
				"two": models.HeaderValue{"b, c"},
			}
			c := Client{BaseHeaders: nil}
			headers := c.makeHeaders(&bh)
			So(headers, ShouldResemble, map[string]string{"one": "a", "two": "b, c"})
		})

		Convey("ClientImpl_makeHeaders with base headers and headers param", func() {
			bh := models.Headers{
				"one": models.HeaderValue{"a"},
				"two": models.HeaderValue{"b, c"},
			}
			hp := models.Headers{
				"three": models.HeaderValue{"d"},
			}
			c := Client{BaseHeaders: &hp}
			headers := c.makeHeaders(&bh)
			So(headers, ShouldResemble, map[string]string{"one": "a", "two": "b, c", "three": "d"})
		})
	})
}

func TestClientImpl_makeRequest(t *testing.T) {
	Convey("ClientImpl_makeRequest", t, func() {
		Convey("ClientImpl_makeRequest succeeds", func() {
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			resp := models.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: body}
			c := Client{
				HttpClient: &mocks.HttpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
			}
			req, _ := http.NewRequest("GET", "https://someurl.com", nil)
			r, err := c.makeRequest(req, c.makeHeaders(&models.Headers{"one": models.HeaderValue{"a"}}))
			So(err, ShouldBeNil)
			So(r, ShouldResemble, resp)
		})

		Convey("ClientImpl_makeRequest fails with fetch error", func() {
			e := errors.New("bad bad thing")
			c := Client{HttpClient: &mocks.HttpClientMock{DoError: e}}
			req, _ := http.NewRequest("GET", "https://someurl.com", nil)
			_, err := c.makeRequest(req, c.makeHeaders(&models.Headers{"one": models.HeaderValue{"a"}}))
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"ClientImpl:makeRequest:failed to make request with error: bad bad thing"}})
		})

		Convey("ClientImpl_makeRequest fails with error reading body", func() {
			e := errors.New("illiterate")
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			c := Client{
				HttpClient: &mocks.HttpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: func(r io.Reader) ([]byte, error) {
					return []byte{}, e
				},
			}
			req, _ := http.NewRequest("GET", "https://someurl.com", nil)
			_, err := c.makeRequest(req, c.makeHeaders(&models.Headers{"one": models.HeaderValue{"a"}}))
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"ClientImpl:makeRequest:failed to read response body with error: illiterate"}})
		})
	})
}

func TestClientImpl_Get(t *testing.T) {
	Convey("ClientImpl_Get", t, func() {
		Convey("ClientImpl_Get succeeds", func() {
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			resp := models.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: body}
			c := Client{
				HttpClient: &mocks.HttpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
				GetRequest: http.NewRequest,
			}
			r, err := c.Get("https://yippie.com", nil, models.UrlParams{models.UrlParam{LVal: "beach", Compare: models.UrlParamTypeEqual, RVal: "ball"}})
			So(err, ShouldBeNil)
			So(r, ShouldResemble, resp)
		})

		Convey("ClientImpl_Get fails when getting new request", func() {
			url := "https://yippie.com"
			e := errors.New("denied")
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			c := Client{
				HttpClient: &mocks.HttpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
				GetRequest: func(method, url string, body io.Reader) (*http.Request, error) {
					return nil, e
				},
			}
			_, err := c.Get(url, nil, models.UrlParams{models.UrlParam{LVal: "beach", Compare: models.UrlParamTypeEqual, RVal: "ball"}})
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"ClientImpl:Get:failed to get new request with url: " + url + "with error: " + e.Error()}})
		})
	})
}

func TestClientImpl_PostBody(t *testing.T) {
	Convey("ClientImpl_PostBody", t, func() {
		Convey("TestClientImpl_PostBody succeeds", func() {
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			resp := models.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: body}
			c := Client{
				HttpClient: &mocks.HttpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
				GetRequest: http.NewRequest,
			}
			r, err := c.PostBody("https://yippie.com", nil, body)
			So(err, ShouldBeNil)
			So(r, ShouldResemble, resp)
		})

		Convey("ClientImpl_PostBody fails when getting new request", func() {
			url := "https://yippie.com"
			e := errors.New("denied")
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			c := Client{
				HttpClient: &mocks.HttpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
				GetRequest: func(method, url string, body io.Reader) (*http.Request, error) {
					return nil, e
				},
			}
			_, err := c.PostBody(url, nil, body)
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"ClientImpl:PostBody:failed to get new request with url: " + url + "with error: " + e.Error()}})
		})
	})
}

func TestClientImpl_PostUrl(t *testing.T) {
	Convey("ClientImpl_PostUrl", t, func() {
		Convey("TestClientImpl_PostUrl succeeds", func() {
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			resp := models.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: body}
			c := Client{
				HttpClient: &mocks.HttpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
				GetRequest: http.NewRequest,
			}
			r, err := c.PostUrl("https://yippie.com", nil, models.UrlParams{models.UrlParam{LVal: "beach", Compare: models.UrlParamTypeEqual, RVal: "ball"}})
			So(err, ShouldBeNil)
			So(r, ShouldResemble, resp)
		})

		Convey("ClientImpl_PostUrl fails when getting new request", func() {
			url := "https://yippie.com"
			e := errors.New("denied")
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			c := Client{
				HttpClient: &mocks.HttpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
				GetRequest: func(method, url string, body io.Reader) (*http.Request, error) {
					return nil, e
				},
			}
			_, err := c.PostUrl(url, nil, models.UrlParams{models.UrlParam{LVal: "beach", Compare: models.UrlParamTypeEqual, RVal: "ball"}})
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"ClientImpl:PostUrl:failed to get new request with url: " + url + "with error: " + e.Error()}})
		})
	})
}
