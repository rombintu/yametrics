package storage

import (
	"errors"
	"strconv"

	"github.com/rombintu/yametrics/internal/metrics"
)

// var ErrNotFound = errors.New("not found")

// type Tables struct {
// 	Counters []m.CounterMetric
// 	Gauges   []m.GaugeMetric
// }

type CounterTable map[string]int64
type GaugeTable map[string]float64

type memDriver struct {
	data map[string]interface{}
}

func NewMemDriver() *memDriver {
	return &memDriver{
		data: make(
			map[string]interface{},
		),
	}
}

func (m *memDriver) Open() error {
	m.data = make(map[string]interface{})
	m.data["counter"] = make(CounterTable)
	m.data["gauge"] = make(GaugeTable)
	return nil
}

func (m *memDriver) Close() error {
	m.data = nil
	return nil
}

func (m *memDriver) GetCounter(key string) (int64, error) {
	data, ok := m.data["counter"].(CounterTable)
	if !ok {
		return 0, errors.New("not found")
	}
	value, exist := data[key]
	if !exist {
		return 0, errors.New("not found")
	}
	return value, nil
}

func (m *memDriver) GetGauge(key string) (float64, error) {
	data, ok := m.data["gauge"].(GaugeTable)
	if !ok {
		return 0, errors.New("not found")
	}
	value, exist := data[key]
	if !exist {
		return 0, errors.New("not found")
	}
	return value, nil
}

func (m *memDriver) UpdateGauge(key string, value float64) {
	data, _ := m.data["gauge"].(GaugeTable)
	data[key] = value
	m.data["gauge"] = data
}

func (m *memDriver) UpdateCounter(key string, value int64) {
	data, _ := m.data["counter"].(CounterTable)
	oldValue := data[key]
	if oldValue == 0 {
		data[key] = value
	} else {
		value = oldValue + value
	}
	data[key] = value
	m.data["counter"] = data
}

func (m *memDriver) GetStorageData() map[string]interface{} {
	return m.data
}

func (m *memDriver) GetMetricByName(mtype, mname string) string {
	switch mtype {
	case metrics.CounterType:
		value, err := m.GetCounter(mname)
		if err != nil {
			return ""
		}
		return strconv.FormatInt(value, 10)
	case metrics.GaugeType:
		value, err := m.GetGauge(mname)
		if err != nil {
			return ""
		}
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return ""
	}
}
