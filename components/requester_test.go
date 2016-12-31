package components

/**
import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
)

type mockRequest struct {
	_Run           func()
	timesRunCalled int
}

func (request *mockRequest) Run() {
	request.timesRunCalled++
	request._Run()
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

func TestDoNothingCauseWereAlreadyDone(t *testing.T) {
	request := &mockRequest{
		_Run: func() {},
	}

	stopper := NewRequester(request)
	stopper.Stop()

	Convey("We shouldn't spin up request because we're canceled immediately", t, func() {
		So(request.timesRunCalled, should.Equal, 0)
		handle, ok := stopper.(*SimpleRequestHandle)
		So(ok, should.BeTrue)
		So(handle.ctx.Err(), should.NotBeNil) // indicates context is done
	})

}

// TODO find good mocking library
func TestSendOneRequestAndShutDown(t *testing.T) {
	stopSignal := make(chan struct{})

	request := &mockRequest{
		_Run: func() {
			stopSignal <- struct{}{}
		},
	}

	handle := NewRequester(request)

	select {
	case <-stopSignal:
		handle.Stop()
	}

	t.Log("wtf")
	Convey("A new requestor should be created, and one request"+
		" should occur before the requestor is stopped", t, func() {
		So(request.timesRunCalled, should.Equal, 1)
		handle, ok := handle.(*SimpleRequestHandle)
		So(ok, should.BeTrue)
		So(handle.ctx.Err(), should.NotBeNil) // indicates context is done
	})
}

func TestRequesterHandle(t *testing.T) {
	handle := newRequestHandle()
	Convey("A request handle should have a cancel function and a context", t, func() {
		So(handle.cancel, should.NotBeNil)
		So(handle.ctx, should.NotBeNil)

		testChan := make(chan int)
		close(testChan) // read the default value in the below select
		var ourChannelIsLegit bool = true

		select {
		case <-handle.ctx.Done():
			ourChannelIsLegit = false
		case <-testChan:
			ourChannelIsLegit = true
		}
		So(ourChannelIsLegit, should.BeTrue)
	})
}
*/
