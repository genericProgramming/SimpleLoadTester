package components

import (
	"testing"

	"fmt"

	"sync/atomic"

	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetNumberToRemove(t *testing.T) {
	var testCases = []struct {
		numCurrentRequestors int
		newNumRequestors     int
		expectedResult       int
	}{
		{0, 0, 0},
		{0, 1, 0},
		{1, 1, 0},
		{1, 0, 1},
		{10, 5, 5},
	}

	Convey("getNumberToRemove should always return a positive number or zero", t, func() {
		for _, tCase := range testCases {
			result := getNumberToRemove(tCase.numCurrentRequestors, tCase.newNumRequestors)
			So(result, should.Equal, tCase.expectedResult)
		}
	})
}

type MockEngine struct {
	_RunAtRate func(r Rate) error
}

func (m *MockEngine) RunAtRate(r Rate) error {
	return m._RunAtRate(r)
}

type MockFactory struct {
	_requestMaker func() (RequestMaker, error)
}

func (factory MockFactory) NewRequestMaker() (RequestMaker, error) {
	return factory._requestMaker()
}

type MockRequestMaker struct {
	_start func() error
	_stop  func() <-chan error
}

func (mrm MockRequestMaker) Start() error {
	return mrm._start()
}

func (mrm MockRequestMaker) Stop() <-chan error {
	return mrm._stop()
}

func TestRunAtRate(t *testing.T) {
	Convey("RequestEngine should error with negative rate", t, func() {
		testEngine := new(RequestEngine)
		negativeValue := Rate(-1)
		err := RateMustNotBeNegative{}
		So(testEngine.RunAtRate(Rate(negativeValue)).Error(), should.Equal, err.Error())
	})

	testCases := []struct {
		stopCalled    bool
		factoryUsed   bool
		rate          Rate
		oldRequestors int
		conveyMessage string
	}{
		{
			stopCalled:    true,
			factoryUsed:   false,
			rate:          Rate(0),
			oldRequestors: 1,
			conveyMessage: "Request Engine should stop all requestors when given a rate of 0",
		},
		{
			stopCalled:    false,
			factoryUsed:   false,
			rate:          Rate(1),
			oldRequestors: 1,
			conveyMessage: "Request Engine should do nothing when the number of requestors is equal to the current number",
		},
		{
			stopCalled:    false,
			factoryUsed:   true,
			rate:          Rate(2),
			oldRequestors: 1,
			conveyMessage: "Request Engine should add requestors when the new rate is greater than the current rate",
		},
	}

	for _, testCase := range testCases {
		Convey(testCase.conveyMessage, t, func() {
			stopCalled := false
			mockMaker := MockRequestMaker{
				func() error { return nil },
				func() <-chan error { stopCalled = true; return nil },
			}
			mockFactory := MockFactory{
				_requestMaker: func() (RequestMaker, error) {
					So(testCase.factoryUsed, should.BeTrue)
					return mockMaker, nil
				},
			}
			testEngine := RequestEngine{
				factory:    mockFactory,
				requestors: []RequestMaker{},
			}
			for i := 0; i < testCase.oldRequestors; i++ {
				testEngine.requestors = append(testEngine.requestors, &mockMaker)
			}

			testEngine.RunAtRate(testCase.rate)
			So(len(testEngine.requestors), should.Equal, int(testCase.rate))
			So(stopCalled, should.Equal, testCase.stopCalled)
		})
	}
}

func TestAddRequestors(t *testing.T) {
	testCases := []struct {
		requestorNumber int
		addNumber       int
		totalNumber     int
	}{
		{0, 1, 1},
		{1, 2, 3},
	}

	for _, testCase := range testCases {
		Convey(fmt.Sprintf("Should increase the requestors by the Rate(%v) with oldNumberRequestors(%v)",
			testCase.addNumber, testCase.requestorNumber), t, func() {
			newRequestCalls := 0
			var startCalled int64 = 0
			mockMaker := MockRequestMaker{
				func() error { atomic.AddInt64(&startCalled, 1); return nil },
				func() <-chan error { So(false, should.BeTrue); return nil },
			}
			mockFactory := MockFactory{
				_requestMaker: func() (RequestMaker, error) {
					newRequestCalls++
					return mockMaker, nil
				},
			}
			requestors := []RequestMaker{}
			for i := testCase.requestorNumber; i > 0; i-- {
				requestors = append(requestors, &mockMaker)
			}
			actualRequestors := addRequestMakers(testCase.addNumber, requestors, mockFactory)
			So(len(actualRequestors), should.Equal, testCase.totalNumber)
			So(newRequestCalls, should.Equal, testCase.addNumber)
			So(startCalled, should.Equal, int64(testCase.addNumber))
		})
	}

}

func TestRemoveRequestors(t *testing.T) {
	testCases := []struct {
		requestorNumber int
		removeNumber    int
	}{
		{1, 1},
		{10, 2},
	}

	Convey("Removing requestors should return a new slice with correct number of makers and call 'stop' on makers removed", t, func() {
		for _, tCase := range testCases {
			stopCallCounter := 0
			mockMaker := MockRequestMaker{
				func() error { return nil },
				func() <-chan error { stopCallCounter++; return nil },
			}
			requestors := []RequestMaker{}
			for i := tCase.requestorNumber; i > 0; i-- {
				requestors = append(requestors, &mockMaker)
			}
			actualRequestors := removeRequestMakers(tCase.removeNumber, requestors)
			So(len(actualRequestors), should.Equal, tCase.requestorNumber-tCase.removeNumber)
			So(stopCallCounter, should.Equal, tCase.removeNumber)
		}
	})
}
