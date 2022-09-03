package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type httpClientMock struct {
	DoResponse http.Response
	DoError    error
}

func (c *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return &c.DoResponse, c.DoError
}

type mockRespObj struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       string `json:"age"`
}

func TestHeaderValue_String(t *testing.T) {
	Convey("HeaderValue_String", t, func() {
		hv := HeaderValue{"one", "two", "three"}
		So(hv.String(), ShouldEqual, "one, two, three")
	})
}

func TestUrlParam_String(t *testing.T) {
	Convey("UrlParam_String", t, func() {
		param := UrlParam{LVal: "beach", Compare: UrlParamTypeEqual, RVal: "ball"}
		So(param.String(), ShouldEqual, "beach=ball")
	})
}

func TestUrlParams_String(t *testing.T) {
	Convey("UrlParams_String", t, func() {
		params := UrlParams{
			UrlParam{LVal: "beach", Compare: UrlParamTypeEqual, RVal: "ball"},
			UrlParam{LVal: "stevie", Compare: UrlParamTypeGreaterThanEqualTo, RVal: "wonder"},
			UrlParam{LVal: "style", Compare: UrlParamTypeLessThan, RVal: "3"},
		}
		So(params.String(), ShouldEqual, "beach=ball?stevie>=wonder?style<3")
	})
}

func TestClientImpl_makeHeaders(t *testing.T) {
	Convey("ClientImpl_makeHeaders", t, func() {
		Convey("ClientImpl_makeHeaders with base headers", func() {
			bh := Headers{
				"one": HeaderValue{"a"},
				"two": HeaderValue{"b, c"},
			}
			c := ClientImpl{BaseHeaders: &bh}
			headers := c.makeHeaders(nil)
			So(headers, ShouldResemble, map[string]string{"one": "a", "two": "b, c"})
		})

		Convey("ClientImpl_makeHeaders with headers param", func() {
			bh := Headers{
				"one": HeaderValue{"a"},
				"two": HeaderValue{"b, c"},
			}
			c := ClientImpl{BaseHeaders: nil}
			headers := c.makeHeaders(&bh)
			So(headers, ShouldResemble, map[string]string{"one": "a", "two": "b, c"})
		})

		Convey("ClientImpl_makeHeaders with base headers and headers param", func() {
			bh := Headers{
				"one": HeaderValue{"a"},
				"two": HeaderValue{"b, c"},
			}
			hp := Headers{
				"three": HeaderValue{"d"},
			}
			c := ClientImpl{BaseHeaders: &hp}
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
			resp := Response{StatusCode: http.StatusOK, Status: "a-ok", Body: body}
			c := ClientImpl{
				HttpClient: &httpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
			}
			req, _ := http.NewRequest("GET", "https://someurl.com", nil)
			r, err := c.makeRequest(req, c.makeHeaders(&Headers{"one": HeaderValue{"a"}}))
			So(err, ShouldBeNil)
			So(r, ShouldResemble, resp)
		})

		Convey("ClientImpl_makeRequest fails with fetch error", func() {
			e := errors.New("bad bad thing")
			c := ClientImpl{HttpClient: &httpClientMock{DoError: e}}
			req, _ := http.NewRequest("GET", "https://someurl.com", nil)
			_, err := c.makeRequest(req, c.makeHeaders(&Headers{"one": HeaderValue{"a"}}))
			So(err, ShouldResemble, e)
		})

		Convey("ClientImpl_makeRequest fails with error reading body", func() {
			e := errors.New("illiterate")
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			c := ClientImpl{
				HttpClient: &httpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: func(r io.Reader) ([]byte, error) {
					return []byte{}, e
				},
			}
			req, _ := http.NewRequest("GET", "https://someurl.com", nil)
			_, err := c.makeRequest(req, c.makeHeaders(&Headers{"one": HeaderValue{"a"}}))
			So(err, ShouldResemble, e)
		})
	})
}

func TestClientImpl_Get(t *testing.T) {
	Convey("ClientImpl_Get", t, func() {
		Convey("ClientImpl_Get succeeds", func() {
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			resp := Response{StatusCode: http.StatusOK, Status: "a-ok", Body: body}
			c := ClientImpl{
				HttpClient: &httpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
			}
			r, err := c.Get("https://yippie.com", nil, UrlParams{UrlParam{LVal: "beach", Compare: UrlParamTypeEqual, RVal: "ball"}})
			So(err, ShouldBeNil)
			So(r, ShouldResemble, resp)
		})
	})
}

func TestClientImpl_PostBody(t *testing.T) {
	Convey("ClientImpl_PostBody", t, func() {
		Convey("TestClientImpl_PostBody succeeds", func() {
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			resp := Response{StatusCode: http.StatusOK, Status: "a-ok", Body: body}
			c := ClientImpl{
				HttpClient: &httpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
			}
			r, err := c.PostBody("https://yippie.com", nil, body)
			So(err, ShouldBeNil)
			So(r, ShouldResemble, resp)
		})
	})
}

func TestClientImpl_PostUrl(t *testing.T) {
	Convey("ClientImpl_PostUrl", t, func() {
		Convey("TestClientImpl_PostUrl succeeds", func() {
			obj := mockRespObj{FirstName: "jimmy", LastName: "fallon", Age: "47"}
			body, _ := json.Marshal(obj)
			resp := Response{StatusCode: http.StatusOK, Status: "a-ok", Body: body}
			c := ClientImpl{
				HttpClient: &httpClientMock{
					DoResponse: http.Response{StatusCode: http.StatusOK, Status: "a-ok", Body: ioutil.NopCloser(bytes.NewReader(body))},
				},
				BodyReader: ioutil.ReadAll,
			}
			r, err := c.PostUrl("https://yippie.com", nil, UrlParams{UrlParam{LVal: "beach", Compare: UrlParamTypeEqual, RVal: "ball"}})
			So(err, ShouldBeNil)
			So(r, ShouldResemble, resp)
		})
	})
}
