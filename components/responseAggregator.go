package components

/**
TODO
Needed metrics are
- Gauge for # of bad requests
- Histogram and ExpSample for response time
- ?? for others
*/

// ResponseAggregator is a Singleton that captures and exposes the metric data
type ResponseAggregator interface {
	Start()
}

type GoMetricBasedAggregator struct {
	config *Config
}

func (aggregator *GoMetricBasedAggregator) Start() {
	// add responseTimeMetrics
	// add status code Gauge
	// TODO figure out what life will look like for
}
