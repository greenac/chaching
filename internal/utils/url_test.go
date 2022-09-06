package utils

import (
	genErr "github.com/greenac/chaching/internal/error"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJoinUrl(t *testing.T) {
	Convey("TestJoinUrl", t, func() {
		Convey("TestJoinUrl should join url when base path has trailing /", func() {
			base := "https://someurl/"
			add := "?value=wee"
			url, err := JoinUrl(base, add)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, base+add)
		})

		Convey("TestJoinUrl should join url when base path does not have trailing / and add begins with /", func() {
			base := "https://someurl"
			add := "/?value=wee"
			url, err := JoinUrl(base, add)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, base+add)
		})

		Convey("TestJoinUrl should join url when base path does not have trailing / and add does not begin with /", func() {
			base := "https://someurl"
			add := "?value=wee"
			url, err := JoinUrl(base, add)
			So(err, ShouldBeNil)
			So(url, ShouldEqual, base+"/"+add)
		})

		Convey("TestJoinUrl should return error when base path is empty", func() {
			_, err := JoinUrl("", "")
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"base url can not be empty string"}})
		})
	})
}
