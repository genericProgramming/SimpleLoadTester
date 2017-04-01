package main

import (
	"net/http"
	"time"

	"fmt"

	"github.com/genericProgramming/simpleLoadTester/components"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	maxRate     = components.Rate(100)
	minRate     = components.Rate(1)
	defaultRate = components.Rate(20)
)

type StairCaseMonitor struct {
	config      *components.Config
	engine      components.Engine
	aggregator  components.ResponseAggregator
	currentRate components.Rate
}

func RunMonitor(monitor *StairCaseMonitor) {
	aggregator := monitor.aggregator.(*components.GoMetricBasedAggregator)
	engine := monitor.engine
	config := monitor.config
	go func() {
		currentRate := monitor.currentRate
		for {
			engine.RunAtRate(currentRate)
			time.Sleep(config.Window)
			histogram := aggregator.ResponseTimeHistogram.Snapshot()
			currentRate = getNewRate(currentRate, histogram, config.Threshold)
			fmt.Println("I've made:", aggregator.CompletedRequests.Count(), "requests and am at ", currentRate)
		}
	}()
}

func getNewRate(currentRate components.Rate, histogram metrics.Histogram, threshold components.ResponseThreshold) components.Rate {
	fmt.Println("Current rate is ", currentRate, "with ave response time", histogram.Percentile(.9), "current threshold", threshold.ResponseTimeMs.Nanoseconds()/1e6)
	if histogram.Percentile(threshold.Level) < float64(threshold.ResponseTimeMs.Nanoseconds()/1e6) {
		newRate := int64(currentRate) * 2
		return components.Rate(min(newRate, int64(maxRate)))
	} else {
		newRate := int64(currentRate) / 2
		return components.Rate(max(newRate, int64(minRate)))
	}
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func createEngine(config *components.Config, resultChannel chan components.RequestResult) components.Engine {
	request := createNewRequest(config, resultChannel)
	requestMaker := components.NewOnePerSecondRequestMakerFactory(request)
	return components.NewRequestEngine(requestMaker)
}

func createNewRequest(config *components.Config, resultChannel chan components.RequestResult) components.Request {
	return components.NewAnnonymousFunctionHttpRequest(func() (*http.Response, error) {
		return http.Get(config.URL)
	}, resultChannel)
}

func NewStairCaseMonitor(config *components.Config) *StairCaseMonitor {
	stairCaseMonitor := &StairCaseMonitor{
		config: config,
	}
	stairCaseMonitor.currentRate = defaultRate
	resultChannel := make(chan components.RequestResult)
	stairCaseMonitor.engine = createEngine(stairCaseMonitor.config, resultChannel)
	stairCaseMonitor.aggregator = components.NewGoMetricBasedAggregator(config)
	stairCaseMonitor.aggregator.ListenAndAggregate(resultChannel)
	return stairCaseMonitor
}
