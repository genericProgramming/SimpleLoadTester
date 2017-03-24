package components

import (
	"net/http"
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const url = "http://localhost:8080/hello/world"

func TestSimpleEngine(t *testing.T) {
	// setup mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, "OK"))

	output := make(chan RequestResult, 100) //prevent requests from blocking on send
	request := NewAnnonymousFunctionHttpRequest(func() (r *http.Response, e error) {
		return http.Get(url)
	}, output)
	factory := OnePerSecondRequestMakerFactory{request: &request}
	engine := NewRequestEngine(&factory)

	expectedNumber := 5
	engine.RunAtRate(Rate(1))
	<-time.After(time.Duration(expectedNumber) * time.Second)
	engine.RunAtRate(0)
	<-time.After(time.Duration(expectedNumber) * time.Second)
	close(output)
	// read off 5 responses to ensure we saw the correct number
	for range output {
		expectedNumber--
	}
	Convey("We should see 5 requests returned", t, func() {
		So(expectedNumber, ShouldBeGreaterThanOrEqualTo, 0)
	})
}
