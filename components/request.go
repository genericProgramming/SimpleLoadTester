package components

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rcrowley/go-metrics"
)

const HttpRequestCounter = "http.metrics.counter"
const HttpRequestTimer = "http.metrics.timer"

type Request interface {
	RunRequest(ctx context.Context)
}

type RequestResult struct {
	StartTime      time.Time
	EndTime        time.Time
	ResponseStatus int
	Error          error
	Assertions     []ResponseAssertionResult //TODO
}

/// Some assertion that
type ResponseAssertionResult interface {
	IsSatisfied() bool
}

type RequestConfig interface {
	MakeHttpCall() (*http.Response, error)
}

type dummyConfig struct {
	httpCall func() (*http.Response, error)
}

func (d *dummyConfig) MakeHttpCall() (*http.Response, error) {
	return d.httpCall()
}

func NewMetricHttpRequest(config RequestConfig, outputChannel chan<- RequestResult) MetricHttpRequest {
	return MetricHttpRequest{
		requestCounter: metrics.GetOrRegisterCounter(HttpRequestCounter, nil),
		requestTimer:   metrics.GetOrRegisterTimer(HttpRequestTimer, nil),
		outputChannel:  outputChannel,
		config:         config,
	}
}

func NewMetricHttpRequestFunctional(httpCall func() (*http.Response, error), outputChannel chan<- RequestResult) MetricHttpRequest {
	return NewMetricHttpRequest(&dummyConfig{httpCall: httpCall}, outputChannel)
}

type MetricHttpRequest struct {
	requestCounter metrics.Counter
	requestTimer   metrics.Timer
	outputChannel  chan<- RequestResult
	config         RequestConfig
}

func (h *MetricHttpRequest) RunRequest(ctx context.Context) {
	h.requestCounter.Inc(1) // TODO determine if this will cause locking issues (perftest this)
	start := time.Now()
	response, e := h.config.MakeHttpCall()
	elapsed := time.Now()

	h.requestTimer.Update(elapsed.Sub(start)) // TODO perftest

	result := RequestResult{
		StartTime: start,
		EndTime:   elapsed,
	}

	if e != nil {
		result.ResponseStatus = -1
		result.Error = e
	} else {
		result.ResponseStatus = response.StatusCode
	}
	h.outputChannel <- result
	log.Println("Request complete with body")
}
