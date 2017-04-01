package components

import (
	"context"
	"net/http"
	"time"
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

type AnnonymousFunctionConfig struct {
	httpCall func() (*http.Response, error)
}

func (d *AnnonymousFunctionConfig) MakeHttpCall() (*http.Response, error) {
	return d.httpCall()
}

func NewHttpRequest(config RequestConfig, outputChannel chan<- RequestResult) HttpRequest {
	return HttpRequest{
		outputChannel: outputChannel,
		config:        config,
	}
}

func NewAnnonymousFunctionHttpRequest(httpCall func() (*http.Response, error), outputChannel chan<- RequestResult) HttpRequest {
	return NewHttpRequest(&AnnonymousFunctionConfig{httpCall: httpCall}, outputChannel)
}

type HttpRequest struct {
	outputChannel chan<- RequestResult
	config        RequestConfig
}

func (h HttpRequest) RunRequest(ctx context.Context) {
	start := time.Now()
	response, e := h.config.MakeHttpCall()
	elapsed := time.Now()

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
}
