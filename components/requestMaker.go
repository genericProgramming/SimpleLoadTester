package components

import (
	"context"
	"log"
	"time"
)

const oNE_SECOND = 1 * time.Second
const rEQUEST_INTERVAL = oNE_SECOND // TODO externalize this?

type RequestMaker interface {
	Start() error
	Stop() <-chan error
}

type RequestMakerFactory interface {
	NewRequestMaker() (RequestMaker, error)
}

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
		for {
			select {
			case <-requestMaker.requestContext.Done():
				log.Println("Stopping request handle") // TODO add key?
				return
			default: // don't block on the done
			}
			ctx, _ := context.WithCancel(requestMaker.requestContext)
			go requestMaker.request.RunRequest(ctx) // TODO add context to this to stop multiple uber long requests?
			requestMaker.limitToOneRequestPerInterval()
		}
	}()
	return nil
}

func (requestMaker *OnePerSecondRequestMaker) limitToOneRequestPerInterval() {
	time.Sleep(rEQUEST_INTERVAL)
}

func (requestMaker *OnePerSecondRequestMaker) Stop() <-chan error {
	requestMaker.cancel()
	errChannel := make(chan error)
	close(errChannel)
	return errChannel
}
