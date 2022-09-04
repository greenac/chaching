package model

import (
	"fmt"
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
		So(rp.Url(), ShouldEqual, url)
	})
}
