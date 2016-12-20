package engine

import (
	"time"
)

const oNE_SECOND = 1 * time.Second
const rEQUEST_INTERVAL = oNE_SECOND // TODO externalize this?

// TODO what does an error look like here? How would the engine know shit's broken?
func NewRequester(requestMaker RequestFactory, requestWriteChannel <-chan RequestResult) Stoppable {
	handle := NewRequestHandle()

	go spinUpNewRequester(handle, requestMaker, requestWriteChannel)

	return &handle
}

func NewRequestHandle() RequesterHandle {
	return RequesterHandle{
		done: make(chan struct{}),
	}
}

type Stoppable interface {
	Stop()
}

type RequesterHandle struct {
	done chan struct{}
}

func (r *RequesterHandle) Done() <-chan struct{} {
	return r.done
}

func (r *RequesterHandle) Stop() {
	close(r.done)
}

func spinUpNewRequester(requesterHandle RequesterHandle, factory RequestFactory, requestWriteChannel <-chan RequestResult) {
	for {
		select {
		case <-requesterHandle.Done():
			// TODO might need to cancel any in-flight requests (req.Run()) by using a context
			return
		}
		req := factory.NewRequest()
		go req.Run(requestWriteChannel)
		limitToOneRequestPerInterval()
	}
}

type RequestFactory interface {
	NewRequest() Request
}

// TODO make this configurable?
func limitToOneRequestPerInterval() {
	time.Sleep(rEQUEST_INTERVAL)
}

type Request interface {
	Run(<-chan RequestResult)
}

type RequestResult struct {
	timeTaken      time.Duration
	responseStatus int
	assertions     []ResponseAssertionResult //TODO
}

// TODO are these errors or some other special interface?
type ResponseAssertionResult interface{}
