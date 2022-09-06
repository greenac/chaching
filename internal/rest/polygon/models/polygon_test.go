package model

import (
	"fmt"
	url2 "net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPolygonAggregateRequestParams_Url(t *testing.T) {
	Convey("TestPolygonAggregateRequestParams_Url", t, func() {
		from := time.Now()
		to := from.Add(30 * time.Second)
		url := fmt.Sprintf("rabbits/range/1/minute/%d/%d?adjusted=false&sort=asc&limit=120", from.UnixMilli(), to.UnixMilli())
		rp := PolygonAggregateRequestParams{
			Name:          "rabbits",
			Multiplier:    1,
			Timespan:      PolygonAggregateTimespanMinute,
			From:          from,
			To:            to,
			SortDirection: PolygonAggregateSortDirectionAsc,
			Limit:         120,
		}
		uri, _ := url2.JoinPath("https://something", rp.Url())
		fmt.Println("url is:", rp.Url(), "full url:", uri)
		So(rp.Url(), ShouldEqual, url)
	})
}
