package components

import (
	"io/ioutil"
	"os"
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const data = `window: 1s
numerOfBadRequestsPerTimeWindow: 100
threshold:
    level: 99
    responseTimeMs: 1ms
requestconfiguration:
    url: "http://somesite:8080/aNeatPath"
    method: POST
    body: '{"key":"value"}'
`

func getAndMakeTempFile(t *testing.T, data string) string {
	file, err := ioutil.TempFile(".", "testFile")
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString(data)
	file.Close()
	return file.Name()
}

func TestLoadBadConfigFile(t *testing.T) {
	fileName := getAndMakeTempFile(t, "not yaml")
	defer os.RemoveAll(fileName) // todo make this cleaner
	Convey("We should get an error when we can't load the yaml file", t, func() {
		_, err := LoadConfig(fileName)
		So(err, ShouldNotBeNil)
	})
}

func TestConfigFile(t *testing.T) {
	fileName := getAndMakeTempFile(t, data)
	defer os.RemoveAll(fileName) // todo make this cleaner

	Convey("The Config should be successfully loaded", t, func() {
		config, err := LoadConfig(fileName)
		So(err, ShouldBeNil)
		So(config.Window, ShouldEqual, 1*time.Second)
		So(config.NumberOfBadRequestsPerTimeWindow, ShouldEqual, 100)
		So(config.Threshold, ShouldNotBeNil)
		So(config.Threshold.Level, ShouldEqual, float64(99))
		So(config.Threshold.ResponseTimeMs, ShouldEqual, 1*time.Millisecond)
		So(config.RequestConfiguration.URL, ShouldEqual, "http://somesite:8080/aNeatPath")
		So(config.RequestConfiguration.Method, ShouldEqual, "POST")
		So(config.RequestConfiguration.Body, ShouldEqual, `{"key":"value"}`)
	})
}
