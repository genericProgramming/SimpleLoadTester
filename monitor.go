package main

import "github.com/genericProgramming/simpleLoadTester/components"

import "time"

type StairCaseMonitor struct {
	config      *components.Config
	engine      components.Engine
	aggregator  components.ResponseAggregator
	currentRate components.Rate
}

func RunMonitor(monitor *StairCaseMonitor) {
	monitor.setup()

	go func() {
		for {
			// do shit
			time.Sleep(monitor.config.Window)
		}
	}()
}

func (stairCaseMonitor *StairCaseMonitor) setup() {

}

func NewStairCaseMonitor(config *components.Config) *StairCaseMonitor {
	return &StairCaseMonitor{
		config: config,
	}
}
