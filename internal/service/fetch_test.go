package fetch

import (
	"fmt"
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/rest/mocks"
	"github.com/greenac/chaching/internal/rest/models"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewFetchService(t *testing.T) {
	Convey("TestNewFetchService", t, func() {
		url := "http://someurl.com"
		c := mocks.ClientMock{}
		fs := NewFetchService(url, &c)
		So(fs, ShouldResemble, &FetchService{Url: url, RestClient: &c})
	})
}

func TestFetchService_Fetch(t *testing.T) {
	Convey("TestFetchService_Fetch", t, func() {
		Convey("TestFetchService_Fetch should succeed", func() {
			url := "http://someurl.com"
			body := []byte("body bytes")
			c := mocks.ClientMock{GetResponse: models.Response{Status: "", StatusCode: http.StatusOK, Body: body}}
			fs := NewFetchService(url, &c)
			b, err := fs.Fetch(models.UrlParams{})
			So(err, ShouldBeNil)
			So(b, ShouldResemble, body)
		})

		Convey("TestFetchService_Fetch should fail with get error", func() {
			ge := &genErr.GenError{Messages: []string{"FetchService:Fetch:failed to get"}}
			c := mocks.ClientMock{GetError: ge}
			fs := NewFetchService("http://someurl.com", &c)
			_, err := fs.Fetch(models.UrlParams{})
			So(err, ShouldResemble, ge)
		})

		Convey("TestFetchService_Fetch should fail with bad request status", func() {
			url := "http://someurl.com"
			body := []byte("body bytes")
			c := mocks.ClientMock{GetResponse: models.Response{Status: "bad", StatusCode: http.StatusBadRequest, Body: body}}
			fs := NewFetchService(url, &c)
			_, err := fs.Fetch(models.UrlParams{})
			So(
				err,
				ShouldResemble,
				&genErr.GenError{Messages: []string{fmt.Sprintf("FetchService:Fetch failed with code: %d and status %s for url: %s", http.StatusBadRequest, "bad", url)}},
			)
		})
	})
}
