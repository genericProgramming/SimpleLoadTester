package components

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregateResults(t *testing.T) {
	Convey("Aggregating Results should always update the elapse time on the histogram", t, func() {
		So(1, should.Equal, 2)
	})

	Convey("Aggregating Results should increment the error code counter on any error returned", t, func() {
		So(1, should.Equal, 2)
	})

	Convey("Aggregating Results should increment total request count for any test", t, func() {
		So(1, should.Equal, 2)
	})
}

func TestAggregateEnd(t *testing.T) {
	Convey("Aggregation should consume forever until the channel passed in is closed", t, func() {
		So(1, should.Equal, 1)
	})
}
