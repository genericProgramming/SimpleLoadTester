package main

import (
	"github.com/genericProgramming/simpleLoadTester/components"
	"github.com/pkg/errors"
)

func main() {
	fileName := components.DefaultFileName // config this
	config, err := components.LoadConfig(fileName)
	if err != nil {
		errors.Wrap(err, "Cant load config file and thus can't start up")
	}
	monitor := NewMonitor(config)
	monitor.Start()
}
