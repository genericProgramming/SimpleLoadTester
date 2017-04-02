package main

import (
	"net/http"
	"time"

	"fmt"

	. "github.com/genericProgramming/simpleLoadTester/components"
)

const (
	maxRate     = Rate(100)
	minRate     = Rate(1)
	defaultRate = Rate(20)
)

var client = http.Client{
	Timeout: 1 * time.Second,
}

type StairCaseMonitor struct {
	config      *Config
	engine      Engine
	aggregator  ResponseAggregator
	currentRate Rate
	isRunning   bool
}

func (monitor *StairCaseMonitor) Stop() {
	monitor.isRunning = false
}

func (monitor *StairCaseMonitor) RunMonitor() {
	monitor.isRunning = true
	aggregator := monitor.aggregator.(*GoMetricBasedAggregator)
	engine := monitor.engine
	config := monitor.config
	go func() {
		currentRate := monitor.currentRate
		for monitor.isRunning {
			engine.RunAtRate(currentRate)
			time.Sleep(config.Window)
			currentRate = getNewRate(currentRate, aggregator, config)
			fmt.Println("I've made:", aggregator.CompletedRequests.Count(), "requests and am at ", currentRate)
		}
	}()
}

func getNewRate(currentRate Rate, aggregator *GoMetricBasedAggregator, config *Config) Rate {
	histogram := aggregator.ResponseTimeHistogram.Snapshot() // TODO the aggregator is way to leaky
	threshold := config.Threshold

	fmt.Println("Current rate is ", currentRate, "with ave response time", histogram.Percentile(.9), "current threshold", threshold.ResponseTimeMs.Nanoseconds()/1e6)

	isUnderThreshold := histogram.Percentile(threshold.Level) < float64(threshold.ResponseTimeMs.Nanoseconds()/1e6)
	tooManyBadRequests := aggregator.ErrorCodesCounter.Count() > config.NumberOfBadRequestsPerTimeWindow

	if !isUnderThreshold || tooManyBadRequests {
		return halfRate(currentRate)
	}
	return doubleRate(currentRate)
}

func doubleRate(currentRate Rate) Rate {
	newRate := int64(currentRate) * 2
	return Rate(min(newRate, int64(maxRate)))
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func halfRate(currentRate Rate) Rate {
	newRate := int64(currentRate) / 2
	return Rate(max(newRate, int64(minRate)))
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func createEngine(config *Config, resultChannel chan RequestResult) Engine {
	request := createNewRequest(config, resultChannel)
	requestMaker := NewOnePerSecondRequestMakerFactory(request)
	return NewRequestEngine(requestMaker)
}

func createNewRequest(config *Config, resultChannel chan RequestResult) Request {
	return NewAnnonymousFunctionHttpRequest(func() (*http.Response, error) {
		httpRequest, err := http.NewRequest(config.RequestConfiguration.Method,
			config.RequestConfiguration.URL, config.RequestConfiguration.Body)
		return client.Do(httpRequest)
	}, resultChannel)
}

func NewStairCaseMonitor(config *Config) *StairCaseMonitor {
	stairCaseMonitor := &StairCaseMonitor{
		config: config,
	}
	stairCaseMonitor.currentRate = defaultRate
	resultChannel := make(chan RequestResult)
	stairCaseMonitor.engine = createEngine(stairCaseMonitor.config, resultChannel)
	stairCaseMonitor.aggregator = NewGoMetricBasedAggregator(config)
	stairCaseMonitor.aggregator.ListenAndAggregate(resultChannel)
	return stairCaseMonitor
}
