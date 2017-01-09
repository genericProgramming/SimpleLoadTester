package components

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"time"

	"github.com/genericProgramming/simpleLoadTester/components/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRunRequest(t *testing.T) {
	testCases := []struct {
		testCode     int
		testError    error
		expectedCode int
	}{
		{1, nil, 1},
		{0, errors.New("Error"), -1},
	}

	for _, tCase := range testCases {
		requestConfig := mocks.RequestConfig{}
		output := make(chan RequestResult, 1)

		req := HttpRequest{
			outputChannel: output,
			config:        &requestConfig,
		}
		requestConfig.On("MakeHttpCall").Return(&http.Response{StatusCode: tCase.testCode}, tCase.testError).After(1 * time.Millisecond)
		ctx := context.Background()
		req.RunRequest(ctx)
		response := <-output
		Convey("Request should leverage requestConfig", t, func() {
			So(response.assertions, ShouldBeNil)
			So(response.err, ShouldEqual, tCase.testError)
			So(response.responseStatus, ShouldEqual, tCase.expectedCode)
			So(response.timeTaken, ShouldBeGreaterThanOrEqualTo, 1*time.Millisecond)
		})
	}
}
