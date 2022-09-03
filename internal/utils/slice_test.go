package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSliceContains(t *testing.T) {
	Convey("TestSliceContains", t, func() {
		sl := []int{1, 2, 3}
		Convey("TestSliceContains should return true when slice has target", func() {
			So(SliceContains(sl, 3), ShouldEqual, true)
		})

		Convey("TestSliceContains should return false when slice does not have target", func() {
			So(SliceContains(sl, 0), ShouldEqual, false)
		})
	})
}
