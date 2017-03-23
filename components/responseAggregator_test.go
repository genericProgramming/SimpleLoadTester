package components

import (
	"testing"

	metrics "github.com/rcrowley/go-metrics"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMetrics(t *testing.T) {
	t.Error("wtf guys")
	Convey("This is a metrics test", t, func() {
		s := metrics.NewUniformSample(1028)
		h := metrics.NewHistogram(s)
		h.Update(10)
		h.Update(20)

		result := h.Percentiles([]float64{.4})
		t.Logf("wtf guys %v", result)
	})
}
