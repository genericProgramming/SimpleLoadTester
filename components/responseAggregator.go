package components

import metrics "github.com/rcrowley/go-metrics"
import "time"

const (
	sampleSize           int     = 1028
	sampleAlpha          float64 = 0.015
	responseTimeKey      string  = "ResponseTimeMetric"
	errorCodeKey         string  = "ErrorCounter"
	completedRequestsKey string  = "CompletedRequestsCounter"
)

/**
TODO
Needed metrics are
- Gauge for # of bad requests
- Histogram and ExpSample for response time
- ?? for others
*/

// ResponseAggregator is a Singleton that captures and exposes the metric data
type ResponseAggregator interface {
	ListenAndAggregate(<-chan RequestResult)
}

type GoMetricBasedAggregator struct {
	config                *Config
	ResponseTimeHistogram metrics.Histogram
	ErrorCodesCounter     metrics.Counter
	CompletedRequests     metrics.Counter
}

func NewGoMetricBasedAggregator(config *Config) *GoMetricBasedAggregator {
	aggregator := GoMetricBasedAggregator{config: config}
	aggregator.ResponseTimeHistogram = getResponseTimeMetricsHistogram(config)
	aggregator.ErrorCodesCounter = getAndRegisterCounter(errorCodeKey)
	aggregator.CompletedRequests = getAndRegisterCounter(completedRequestsKey)
	return &aggregator
}

func getResponseTimeMetricsHistogram(config *Config) metrics.Histogram {
	s := metrics.NewExpDecaySample(sampleSize, sampleAlpha) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	metrics.Register(responseTimeKey, h)
	return h
}

func getAndRegisterCounter(key string) metrics.Counter {
	g := metrics.NewCounter()
	metrics.Register(key, g)
	return g
}

// ListenAndAggregate sets up metrics and listens on the channel for completed requests
// aggregating accordingly
// TODO create a stop method on this?
func (aggregator *GoMetricBasedAggregator) ListenAndAggregate(results <-chan RequestResult) {
	go aggregator.aggregateResults(results)
}

func (aggregator *GoMetricBasedAggregator) aggregateResults(results <-chan RequestResult) {
	for result := range results {
		timeTaken := result.EndTime.Sub(result.StartTime) * time.Millisecond
		aggregator.ResponseTimeHistogram.Update(int64(timeTaken.Seconds() / 1000))

		if !requestWasSuccessful(&result) {
			aggregator.ErrorCodesCounter.Inc(1)
		}
		aggregator.CompletedRequests.Inc(1)
	}
}

func requestWasSuccessful(request *RequestResult) bool {
	return request.Error != nil &&
		request.ResponseStatus < 500
}
