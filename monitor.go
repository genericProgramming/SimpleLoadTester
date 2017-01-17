package main

import "github.com/genericProgramming/simpleLoadTester/components"
import "log"

type Monitor interface {
	Start()
}

type StairCaseMonitor struct {
	config *components.Config
	engine components.Engine
}

func (monitor *StairCaseMonitor) Start() {
	log.Panicln("Lol -- probs should implement this")
}

func NewMonitor(config *components.Config) Monitor {
	return &StairCaseMonitor{
		config: config,
	}
}
