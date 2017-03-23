package main

import "github.com/genericProgramming/simpleLoadTester/components"
import "log"
import "time"

type Monitor interface {
	IsApiHealthy()
}

type StairCaseMonitor struct {
	config      *components.Config
	engine      components.Engine
	aggregator  components.ResponseAggregator
	currentRate components.Rate
}

func RunMonitor(monitor *StairCaseMonitor) {
	go func() {
		for {
			// getCurrentStats()
			// compareStatsToSpecification()
			// adjustEngine()
			time.Sleep(monitor.config.Window)
		}
	}()
}

func (monitor *StairCaseMonitor) IsApiHealthy() {
	log.Panic("NYI")
}

func NewMonitor(config *components.Config) Monitor {
	return &StairCaseMonitor{
		config: config,
	}
}
