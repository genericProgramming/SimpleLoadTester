package components

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"

	. "github.com/stretchr/testify/mock"
)

type MockRequest struct {
	_Run           func(ctx context.Context)
	timesRunCalled int
}

func (request *MockRequest) RunRequest(ctx context.Context) {
	request.timesRunCalled++
	request._Run(ctx)
}

func TestLongRequests(t *testing.T) {
	Convey("We should not block when making requests", t, func() {
		stopSignal := make(chan struct{})
		var requestCounter int
		LONG_INTERVAL := time.Second * 2

		request := &mockRequest{
			_Run: func() {
				time.Sleep(LONG_INTERVAL)
				requestCounter++ // sketchy ik
				if requestCounter > 3 {
					close(stopSignal)
				}
			},
		}

		stopper := NewRequester(request)
		select {
		case <-stopSignal:
			stopper.Stop()
		}

		So(requestCounter, should.BeGreaterThanOrEqualTo, 3) // TODO observe context cancel?
		So(request.timesRunCalled, should.BeGreaterThan, requestCounter)
		handle, ok := stopper.(*SimpleRequestHandle)
		So(ok, should.BeTrue)
		So(handle.ctx.Err(), should.NotBeNil) // indicates context is done
	})
}
