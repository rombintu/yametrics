package metrics

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

// type Metrics struct {
// 	Counters []CounterMetric
// 	Gauges   []GaugeMetric
// }
