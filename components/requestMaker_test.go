package components

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
)

type mockRequest struct {
	_Run           func(ctx context.Context)
	timesRunCalled int
}

func (request *mockRequest) RunRequest(ctx context.Context) {
	request.timesRunCalled++
	request._Run(ctx)
}

func TestLongRequests(t *testing.T) {
	Convey("We should not block when making requests", t, func() {
		countSignal := make(chan struct{})
		LONG_INTERVAL := time.Second * 2

		request := &mockRequest{
			_Run: func(ctx context.Context) {
				time.Sleep(LONG_INTERVAL)
				countSignal <- struct{}{}
			},
		}

		stopper := newOnePerSecondRequestMaker(request)
		stopper.Start()

		stopCounter := 0
		for {
			<-countSignal
			stopCounter++
			if stopCounter > 3 {
				break
			}
		}
		stopper.Stop()

		So(request.timesRunCalled, should.BeGreaterThan, stopCounter)
		So(stopper.requestContext.Err(), should.NotBeNil) // indicates context is done
	})
}

func TestStart(t *testing.T) {

}
