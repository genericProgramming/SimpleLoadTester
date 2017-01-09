package components

import (
	"fmt"
	"net/http"
	"testing"

	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const url = "http://localhost:8080/hello/world"

func TestSimpleEngine(t *testing.T) {
	// setup mockys
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, "OK"))

	output := make(chan RequestResult, 100) //prevent requests from blocking on send
	request := NewHttpRequestFunctional(func() (r *http.Response, e error) {
		return http.Get(url)
	}, output)
	factory := OnePerSecondRequestMakerFactory{request: &request}
	engine := NewRequestEngine(&factory)
	engine.RunAtRate(Rate(1))
	<-time.After(5 * time.Second)
	engine.RunAtRate(0)
	<-time.After(5 * time.Second)
	close(output)
	// read off 10 response
	count := 0
	for range output {
		count++
	}
	fmt.Println(count)
}
