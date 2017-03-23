package components

import (
	"context"
	"log"
	"math/rand"
	"time"
)

const oNE_SECOND = 1 * time.Second
const oNE_SECOND_MILLIS = int32(1 * time.Millisecond)
const rEQUEST_INTERVAL = oNE_SECOND // TODO externalize this?

func init() {
	rand.Seed(time.Now().Unix())
}

type RequestMaker interface {
	Start() error
	Stop() <-chan error
}

type RequestMakerFactory interface {
	NewRequestMaker() (RequestMaker, error)
}

// TODO create a NewOnePerSecondRequestMaker method for this
type OnePerSecondRequestMakerFactory struct {
	request Request
}

func (factory *OnePerSecondRequestMakerFactory) NewRequestMaker() (RequestMaker, error) {
	if factory == nil || factory.request == nil {
		log.Panicln("Need to have a valid request struct")
	}
	return NewOnePerSecondRequestMaker(factory.request), nil
}

func NewOnePerSecondRequestMaker(request Request) *OnePerSecondRequestMaker {
	requestMaker := OnePerSecondRequestMaker{}
	requestMaker.request = request
	ctx, cancel := context.WithCancel(context.Background())
	requestMaker.requestContext = ctx
	requestMaker.cancel = cancel
	return &requestMaker
}

type OnePerSecondRequestMaker struct {
	request        Request
	requestContext context.Context
	cancel         context.CancelFunc
}

func (requestMaker *OnePerSecondRequestMaker) Start() error {
	go func() {
		randomlySpaceRequestMakers()
		for {
			select {
			case <-requestMaker.requestContext.Done():
				log.Println("Stopping request handle") // TODO add key?
				return
			default:
			}
			ctx := context.WithValue(requestMaker.requestContext, struct{}{}, struct{}{})
			go requestMaker.request.RunRequest(ctx) // TODO who decides if this is run in a separate go routine?
			limitToOneRequestPerInterval()
		}
	}()
	return nil
}

// TODO this probably belongs in the engine somewhere
func randomlySpaceRequestMakers() {
	sleepTime := time.Duration(rand.Int31n(1000)) * time.Millisecond
	time.Sleep(sleepTime)
}

func limitToOneRequestPerInterval() {
	time.Sleep(rEQUEST_INTERVAL)
}

// Stop prevents new requests from being created and signals the inflight requests that the should stop
func (requestMaker *OnePerSecondRequestMaker) Stop() <-chan error {
	requestMaker.cancel()
	errChannel := make(chan error)
	close(errChannel)
	return errChannel
}
