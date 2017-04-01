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
	RequestConfiguration             RequestConfiguration
}

type ResponseThreshold struct {
	Level          float64       // TODO is this an enum?
	ResponseTimeMs time.Duration `yaml:"responseTimeMs"`
}

type RequestConfiguration struct {
	URL    string
	Method string
	Body   string
}

func LoadConfig(fileName string) (*Config, error) {
	config := new(Config)
	configData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read input configuration file")
	}
	err = yaml.Unmarshal(configData, &config)
	return config, err
}
