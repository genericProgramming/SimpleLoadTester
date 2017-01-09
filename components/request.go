package components

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Request interface {
	RunRequest(ctx context.Context)
}

type RequestResult struct {
	timeTaken      time.Duration
	responseStatus int
	err            error
	assertions     []ResponseAssertionResult //TODO
}

type ResponseAssertionResult interface {
	IsSatisfied() bool
}

type RequestConfig interface {
	MakeHttpCall() (*http.Response, error)
}

func NewHttpRequest(config RequestConfig, outputChannel chan<- RequestResult) HttpRequest {
	return HttpRequest{
		outputChannel: outputChannel,
		config:        config,
	}
}

func NewHttpRequestFunctional(httpCall func() (*http.Response, error), outputChannel chan<- RequestResult) HttpRequest {
	return NewHttpRequest(&dummyConfig{httpCall: httpCall}, outputChannel)
}

type dummyConfig struct {
	httpCall func() (*http.Response, error)
}

func (d *dummyConfig) MakeHttpCall() (*http.Response, error) {
	return d.httpCall()
}

type HttpRequest struct {
	outputChannel chan<- RequestResult
	config        RequestConfig
}

// TODO use the context in the request
func (h *HttpRequest) RunRequest(ctx context.Context) {
	// simulate work
	start := time.Now()
	response, e := h.config.MakeHttpCall()
	elapsed := time.Since(start)

	result := RequestResult{
		timeTaken: elapsed,
	}

	if e != nil {
		result.responseStatus = -1
		result.err = e
	} else {
		result.responseStatus = response.StatusCode
	}
	h.outputChannel <- result
	log.Println("Request complete with body")
}
