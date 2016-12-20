package engine

import (
	"testing"
)

func TestDoNothingBecauseWereDone(t *testing.T){
	mockHandle := RequesterHandle{}
}

/*
func TestDoNothingBecauseImDone(t *testing.T){
	done, cancel := context.WithCancel(context.Background())
	result := make(chan RequestResult)
	requestSynthesizer := MockRequestSynthesizer{wg: sync.WaitGroup{}, numberOfRunCalls:-10, numberOfGetARequestCalls:-9}

	Convey("No Requests should be called", t, func(){
		cancel()
		NewRequester(done, result, &requestSynthesizer)
		So(requestSynthesizer.numberOfRunCalls, should.Equal, -10)
		So(requestSynthesizer.numberOfGetARequestCalls, should.Equal, -9)
	})
}

func TestCreateNewRequestAndShutDown(t *testing.T) {
	t.Skip(123)
	done, cancel := context.WithCancel(context.Background())
	result := make(chan RequestResult,100)
	requestSynthesizer := MockRequestSynthesizer{wg: sync.WaitGroup{}}

	err := NewRequester(done, result, requestSynthesizer)

	if err != nil {
		t.Errorf("Couldn't create a requester %v", err)
	}

	time.Sleep(2 * time.Second)
	cancel()

	Convey("The Requester should've run a few times", t, func() {
		So(requestSynthesizer.numberOfRunCalls, should.BeGreaterThan, 0)
	})
}
*/
