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
	requestors []RequestMaker
	factory    RequestMakerFactory
}

func (engine *RequestEngine) RunAtRate(rate Rate) error {
	newNumberOfRequestMakers := int(rate)
	return engine.updateRequestMaker(newNumberOfRequestMakers)
}

// TODO there's a cleaner way to do this -- figure it out
func (engine *RequestEngine) updateRequestMaker(newNumRequestMaker int) error {
	if newNumRequestMaker < 0 {
		return RateMustNotBeNegative{}
	}

	numCurrentRequestMaker := len(engine.requestors)
	var newRequestMakers []RequestMaker
	if newNumRequestMaker > numCurrentRequestMaker {
		numberToAdd := newNumRequestMaker - numCurrentRequestMaker
		newRequestMakers = addRequestMakers(numberToAdd, engine.requestors, engine.factory)
	} else {
		numberToRemove := getNumberToRemove(numCurrentRequestMaker, newNumRequestMaker)
		newRequestMakers = removeRequestMakers(numberToRemove, engine.requestors)
	}

	engine.requestors = newRequestMakers

	return nil
}

type RateMustNotBeNegative struct{}

func (e RateMustNotBeNegative) Error() string {
	return "Rate must not be negative"
}

func addRequestMakers(howManyToAdd int, requestors []RequestMaker, factory RequestMakerFactory) []RequestMaker {
	lenRequestMaker := len(requestors)
	totalNewRequestMaker := lenRequestMaker + howManyToAdd

	newRequestMaker := make([]RequestMaker, totalNewRequestMaker)
	copy(newRequestMaker, requestors)

	for i := lenRequestMaker; i < totalNewRequestMaker; i++ {
		newRequestMaker[i], _ = factory.NewRequestMaker() // TODO use this
	}
	return newRequestMaker
}

func getNumberToRemove(currentRequestMakers int, newRequestMakers int) int {
	numberToRemove := currentRequestMakers - newRequestMakers
	if numberToRemove < 0 {
		numberToRemove = currentRequestMakers
	}
	return numberToRemove
}

func removeRequestMakers(numberToRemove int, requestors []RequestMaker) []RequestMaker {
	lenRequestMaker := len(requestors)
	requestMakersToStop := requestors[lenRequestMaker-numberToRemove : lenRequestMaker]
	for _, requestMaker := range requestMakersToStop {
		requestMaker.Stop()
	}
	return requestors[:lenRequestMaker-numberToRemove]
}
