package components

import (
	"context"
	"time"
)

type Request interface {
	RunRequest(ctx context.Context)
}

type RequestResult struct {
	timeTaken      time.Duration
	responseStatus int
	assertions     []ResponseAssertionResult //TODO
}

type ResponseAssertionResult interface {
	IsSatisfied() bool
}

func NewHttpRequest(outputChannel chan<- RequestResult) HttpRequest {
	return HttpRequest{
		outputChannel,
	}
}

type HttpRequest struct {
	outputChannel chan<- RequestResult
}

// TODO
func (h *HttpRequest) RunRequest(ctx context.Context) {
	// simulate work
	time.Sleep(time.Millisecond * 10)
}
