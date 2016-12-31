package components

/**
import (
	"context"
	"fmt"
	"time"
)

const oNE_SECOND = 1 * time.Second
const rEQUEST_INTERVAL = oNE_SECOND // TODO externalize this?

// Creates a new request
type RequestFactory interface {
	NewRequest() Request
}

type Request interface {
	Run(chan<- RequestResult)
}

type OnePerSecond struct{}

func (onePerSecond *OnePerSecond) Run() RequestHandle {
	handle := newRequestHandle()

	go startRequests(handle)

	return &handle
}

func newRequestHandle() SimpleRequestHandle {
	ctx, cancel := context.WithCancel(context.Background())
	return SimpleRequestHandle{
		ctx:    ctx,
		cancel: cancel,
	}
}

type RequestHandle interface {
	Stop()
}

// TODO make this an interface
type SimpleRequestHandle struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (r *SimpleRequestHandle) Stop() {
	r.cancel()
}

func (r *SimpleRequestHandle) startRequests(request Request) {
	for {
		select {
		case <-r.ctx.Done():
			fmt.Println("Stopping request handle") // TODO add key?
			return
		default: // don't block on the done
		}
		go request.Run() // TODO add context to this to stop multiple uber long requests?
		limitToOneRequestPerInterval()
	}
}

// TODO make this configurable?
func limitToOneRequestPerInterval() {
	time.Sleep(rEQUEST_INTERVAL)
}

type RequestResult struct {
	timeTaken      time.Duration
	responseStatus int
	assertions     []ResponseAssertionResult //TODO
}

// TODO are these errors or some other special interface?
type ResponseAssertionResult interface{}
*/
