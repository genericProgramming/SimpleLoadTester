package components

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestGetNumberToRemove(t *testing.T) {
	var testCases = []struct {
		numCurrentRequestors int
		newNumRequestors     int
		expectedResult       int
	}{
		{0, 0, 0},
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
	mock.Mock
	RequestEngine
}

func (m *MockEngine) RunAtRate(r Rate) error {
	args := m.Called(r)
	return args.Error(0)
}

func TestRunAtRate(t *testing.T) {
	Convey("RequestEngine should error with negative rate", t, func() {
		testEngine := new(RequestEngine)
		negativeValue := Rate(-1)
		err := RateMustNotBeNegative{}
		So(testEngine.RunAtRate(Rate(negativeValue)).Error(), should.Equal, err.Error())
	})

	Convey("Request Engine should be cool with a non-negative rate", t, func() {
		rates := []Rate{
			Rate(0),
			Rate(1),
		}
		testEngine := new(MockEngine)
		// var testFunc func(*HttpRequestEngine, Rate) error
		// testFunc = (*HttpRequestEngine).RunAtRate
		for _, rate := range rates {
			testEngine.On("updateRequestors", int(rate)).Return()
			//So(testFunc(testEngine.(HttpRequestEngine), rate), should.BeNil)
		}
	})
}

type MockFactory struct {
	mock.Mock
}

func (factory *MockFactory) NewRequest() RequestHandle {
	args := factory.Called()
	return args.Get(0).(RequestHandle)
}

type MockRequestHandle struct{}

func (e *MockRequestHandle) Stop() {}

func TestAddRequestors(t *testing.T) {
	testCases := []struct {
		howManyToAdd       int
		requestors         []RequestHandle
		request            RequestHandle
		expectedRequestors []RequestHandle
	}{
		{
			1,
			make([]RequestHandle, 1),
			MockRequestHandle{},
			[]RequestHandle{},
		},
	}
}
