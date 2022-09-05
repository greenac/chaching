package fetch

import (
	"errors"
	"fmt"
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/rest/mocks"
	"github.com/greenac/chaching/internal/rest/models"
	polygonModels "github.com/greenac/chaching/internal/rest/polygon/models"
	"net/http"
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewFetchService(t *testing.T) {
	Convey("TestNewFetchService", t, func() {
		uri := "http://someurl.com"
		c := mocks.ClientMock{}
		fs := NewFetchService(uri, &c, url.JoinPath)
		So(fs, ShouldResemble, &FetchService{Url: uri, RestClient: &c})
	})
}

func TestFetchService_Fetch(t *testing.T) {
	Convey("TestFetchService_Fetch", t, func() {
		Convey("TestFetchService_Fetch should succeed", func() {
			uri := "http://someurl.com"
			body := []byte("body bytes")
			c := mocks.ClientMock{GetResponse: models.Response{Status: "", StatusCode: http.StatusOK, Body: body}}
			fs := NewFetchService(uri, &c, url.JoinPath)
			b, err := fs.Fetch(models.UrlParams{})
			So(err, ShouldBeNil)
			So(b, ShouldResemble, body)
		})

		Convey("TestFetchService_Fetch should fail with get error", func() {
			ge := &genErr.GenError{Messages: []string{"FetchService:Fetch:failed to get"}}
			c := mocks.ClientMock{GetError: ge}
			fs := NewFetchService("http://someurl.com", &c, url.JoinPath)
			_, err := fs.Fetch(models.UrlParams{})
			So(err, ShouldResemble, ge)
		})
	})
}

func TestFetchService_FetchWithFetchData(t *testing.T) {
	Convey("TestFetchService_FetchWithFetchData", t, func() {
		from := time.Now()
		to := from.Add(30 * time.Second)
		fd := polygonModels.PolygonAggregateRequestParams{
			Name:          "rabbits",
			Multiplier:    1,
			Timespan:      polygonModels.PolygonAggregateTimespanMinute,
			From:          from,
			To:            to,
			SortDirection: polygonModels.PolygonAggregateSortDirectionAsc,
			Limit:         120,
		}

		Convey("TestFetchService_FetchWithFetchData should succeed", func() {
			uri := "http://someurl.com"
			body := []byte("body bytes")
			c := mocks.ClientMock{GetResponse: models.Response{Status: "", StatusCode: http.StatusOK, Body: body}}
			fs := NewFetchService(uri, &c, url.JoinPath)
			b, err := fs.FetchWithFetchData(fd)
			So(err, ShouldBeNil)
			So(b, ShouldResemble, body)
		})

		Convey("TestFetchService_FetchWithFetchData should fail when joining url", func() {
			e := errors.New("bad bad thing")
			uri := "http://someurl.com"
			expUrl := fmt.Sprintf("rabbits/range/1/minute/%d/%d?adjusted=false&sort=asc&limit=120", from.UnixMilli(), to.UnixMilli())
			body := []byte("body bytes")
			c := mocks.ClientMock{GetResponse: models.Response{Status: "", StatusCode: http.StatusOK, Body: body}}
			fs := NewFetchService(uri, &c, func(base string, elem ...string) (result string, err error) {
				return "", e
			})
			_, err := fs.FetchWithFetchData(fd)
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"FetchService:FetchWithFetchData failed to join urls: http://someurl.com and " + expUrl}})
		})

		Convey("TestFetchService_FetchWithFetchData should fail when client gets", func() {
			e := errors.New("bad bad thing")
			uri := "http://someurl.com"
			c := mocks.ClientMock{GetError: &genErr.GenError{Messages: []string{e.Error()}}}
			fs := NewFetchService(uri, &c, url.JoinPath)
			_, err := fs.FetchWithFetchData(fd)
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{e.Error(), "FetchService:Fetch:failed to get"}})
		})
	})
}
func TestFetchService_handleResponse(t *testing.T) {
	Convey("TestFetchService_handleResponse", t, func() {
		Convey("TestFetchService_handleResponse should succeed", func() {
			uri := "http://someurl.com"
			body := []byte("body bytes")
			fs := NewFetchService(uri, nil, url.JoinPath)
			b, err := fs.handleResponse(models.Response{Status: "ok", StatusCode: http.StatusOK, Body: body})
			So(err, ShouldBeNil)
			So(b, ShouldResemble, body)
		})

		Convey("TestFetchService_handleResponse should fail with bad request status", func() {
			uri := "http://someurl.com"
			fs := NewFetchService(uri, nil, url.JoinPath)
			_, err := fs.handleResponse(models.Response{StatusCode: 400, Status: "bad"})
			So(
				err,
				ShouldResemble,
				&genErr.GenError{Messages: []string{fmt.Sprintf("FetchService:handleResponse failed with code: %d and status %s for url: %s", http.StatusBadRequest, "bad", uri)}},
			)
		})
	})
}
