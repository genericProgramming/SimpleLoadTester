package components

type Engine interface {
	Throttle
}

type Throttle interface {
	RunAtRate(Rate) error
}

type Rate int

// TODO add stats to this?
type RequestEngine struct {
	requestors []RequestHandle
	factory    RequestFactory
}

func (engine *RequestEngine) RunAtRate(rate Rate) error {
	newNumRequestors := int(rate)
	return engine.updateRequestors(newNumRequestors)
}

type RateMustNotBeNegative struct{}

func (e RateMustNotBeNegative) Error() string {
	return "Rate must not be negative"
}

// TODO there's a cleaner way to do this -- figure it out
func (engine *RequestEngine) updateRequestors(newNumRequestors int) error {
	if newNumRequestors < 0 {
		return RateMustNotBeNegative{}
	}

	numCurrentRequestors := len(engine.requestors)
	var newRequestors []RequestHandle
	if newNumRequestors > numCurrentRequestors {
		numberToAdd := newNumRequestors - numCurrentRequestors
		newRequestors = addRequestors(numberToAdd, engine.requestors, engine.factory)
	} else {
		numberToRemove := getNumberToRemove(newNumRequestors, numCurrentRequestors)
		newRequestors = removeRequestors(numberToRemove, engine.requestors)
	}
	engine.requestors = newRequestors

	return nil
}

func addRequestors(howManyToAdd int, requestors []RequestHandle, factory RequestFactory) []RequestHandle {
	lenRequestors := len(requestors)
	totalNewRequestors := lenRequestors + howManyToAdd

	newRequestors := make([]RequestHandle, lenRequestors, totalNewRequestors)
	copy(newRequestors, requestors)

	for i := lenRequestors; i < totalNewRequestors; i++ {
		newRequestors[i] = factory.NewRequest()
	}
	return newRequestors
}

func getNumberToRemove(numCurrentRequestors int, newNumRequestors int) int {
	numberToRemove := numCurrentRequestors - newNumRequestors
	if numberToRemove < 0 {
		numberToRemove = numCurrentRequestors
	}
	return numberToRemove
}

func removeRequestors(numberRemove int, requestors []RequestHandle) []RequestHandle {
	lenRequestors := len(requestors)
	toBeRemovedAndStopped := requestors[lenRequestors-numberRemove : lenRequestors]
	for _, stoppable := range toBeRemovedAndStopped {
		stoppable.Stop()
	}
	return requestors[:lenRequestors-numberRemove]
}
