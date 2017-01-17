package main

import "github.com/genericProgramming/simpleLoadTester/components"
import "log"

type Monitor interface {
	IsApiHealthy()
	RampUp()
	RampDown()
}

type StairCaseMonitor struct {
	config      *components.Config
	engine      components.Engine
	aggregator  components.ResponseAggregator
	currentRate components.Rate
}

func RunMonitor(monitor Monitor) {
	go func(){
		for {
			performWork()
			sleep()
		}
	}()
}

func isApiHealthy(aggregator components.ResponseAggregator) Data{
	log.Panic("NYI")
}

func 

func NewMonitor(config *components.Config) Monitor {
	return &StairCaseMonitor{
		config: config,
	}
}
