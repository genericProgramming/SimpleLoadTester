package components

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const data = `
window: 1s
numerOfBadRequestsPerTimeWindow: 100
threshold:
    level: 99
    responseTimeMs: 1ms`

func getAndMakeTempFile(t *testing.T) string {
	file, err := ioutil.TempFile(".", "testFile")
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString(data)
	file.Close()
	return file.Name()
}

func TestConfigFile(t *testing.T) {
	fileName := getAndMakeTempFile(t)
	defer os.RemoveAll(fileName) // todo make this cleaner

	fmt.Println(time.ParseDuration(" 1s"))

	Convey("The Config should be successfully loaded", t, func() {
		config, err := LoadConfig(fileName)
		fmt.Println("---", config)
		So(err, ShouldBeNil)
		So(config.Window, ShouldEqual, 1*time.Second)
		So(config.NumberOfBadRequestsPerTimeWindow, ShouldEqual, 100)
		So(config.Threshold, ShouldNotBeNil)
		So(config.Threshold.Level, ShouldEqual, float64(99))
		So(config.Threshold.ResponseTimeMs, ShouldEqual, 1*time.Millisecond)
	})
}
