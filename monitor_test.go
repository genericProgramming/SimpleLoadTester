package main

import (
	"testing"

	"time"

	. "github.com/genericProgramming/simpleLoadTester/components"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRateCalculation(t *testing.T) {
	currentRate := Rate(10)
	config := &Config{
		NumberOfBadRequestsPerTimeWindow: 10,
		Threshold: ResponseThreshold{
			Level:          .9,
			ResponseTimeMs: 100 * time.Millisecond,
		},
	}
	aggregator := NewGoMetricBasedAggregator(config)

	Convey("The rate should be cut in half when the number of bad requests passes the threshold", t, func() {
		aggregator.ErrorCodesCounter.Inc(11)
		aggregator.ResponseTimeHistogram.Update(1)
		newRate := getNewRate(currentRate, aggregator, config)
		So(newRate, should.Equal, Rate(5))
	})

	Convey("The rate should be cut in half when the thresold response times drops", t, func() {
		aggregator.ErrorCodesCounter.Clear()
		aggregator.ResponseTimeHistogram.Clear()
		aggregator.ResponseTimeHistogram.Update(101)
		newRate := getNewRate(currentRate, aggregator, config)
		So(newRate, should.Equal, Rate(5))
	})

	Convey("The rate should double when all is good", t, func() {
		aggregator.ErrorCodesCounter.Clear()
		aggregator.ResponseTimeHistogram.Clear()
		newRate := getNewRate(currentRate, aggregator, config)
		So(newRate, should.Equal, Rate(20))
	})
}

func TestMetrics(t *testing.T) {
	Convey("Hammer that service", t, func() {
		// config := &Config{
		// 	NumberOfBadRequestsPerTimeWindow: 10,
		// 	Threshold: ResponseThreshold{
		// 		Level:          .9,
		// 		ResponseTimeMs: 100 * time.Millisecond,
		// 	},
		// 	URL:    "http://localhost:8081/hello/world",
		// 	Window: 10 * time.Second,
		// }
		// monitor := NewStairCaseMonitor(config)
		// monitor.RunMonitor()
		// <-time.After(100 * time.Second)
		t.Log("We're done team!")
	})
}
