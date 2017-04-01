package components

import (
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const DefaultFileName = "TestConfig.yaml"

type Config struct {
	Window                           time.Duration
	NumberOfBadRequestsPerTimeWindow int64 `yaml:"numerOfBadRequestsPerTimeWindow"`
	Threshold                        ResponseThreshold
	URL                              string
}

type ResponseThreshold struct {
	Level          float64       // TODO is this an enum?
	ResponseTimeMs time.Duration `yaml:"responseTimeMs"`
}

func LoadConfig(fileName string) (*Config, error) {
	config := new(Config)
	configData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read input configuration file")
	}
	yaml.Unmarshal(configData, &config)
	return config, nil
}
