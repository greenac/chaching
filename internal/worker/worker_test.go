package worker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWorker(t *testing.T) {
	Convey("TestWorker_Work", t, func() {
		expResult1 := "asip"
		expResult2 := "bandits"

		c := make(chan Message[string], 2)
		w := NewWorker[string](5, c)
		w.Work()

		w.AddTask(func() string {
			return expResult1
		})

		w.AddTask(func() string {
			return expResult2
		})

		m1 := <-c
		m2 := <-c

		results := map[string]bool{m1.Result: true, m2.Result: true}
		So(results, ShouldContainKey, expResult1)
		So(results, ShouldContainKey, expResult2)
	})
}
