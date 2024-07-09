package storage

import (
	"errors"
	"strconv"

	"github.com/rombintu/yametrics/internal/metrics"
)

var ErrInvalidDriver = errors.New("invalid storage driver. Available drivers: memory")

const (
	memDriverType  = "memory"
	mockDriverType = "mock"
	// TODO: add more drivers
	// fileDriverType  = "file"
	// redisDriverType = "redis"
	// mongoDriverType = "mongo"
)

type StorageDriver interface {
	Open() error
	Close() error
	GetCounter(key string) (int64, error)
	GetGauge(key string) (float64, error)
	UpdateGauge(key string, value float64)
	UpdateCounter(key string, value int64)
	GetStorageData() map[string]interface{}
	GetMetricByName(mtype, mname string) string
}

type Storage struct {
	Driver StorageDriver
}

func NewStorage(driverType string) (*Storage, error) {
	var driver StorageDriver
	switch driverType {
	case memDriverType, mockDriverType:
		driver = NewMemDriver()
	default:
		return nil, ErrInvalidDriver
	}
	return &Storage{
		Driver: driver,
	}, nil
}

func (s *Storage) Open() error {
	return s.Driver.Open()
}

func (s *Storage) Close() error {
	return s.Driver.Close()
}

func (s *Storage) GetAllMetrics() map[string]interface{} {
	return s.Driver.GetStorageData()
}

func (s *Storage) WriteMetric(metricType, metricName string, metricValue string) error {
	switch metricType {
	case metrics.GaugeType:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("cannot parse float64 from gauge metric")
		}

		s.Driver.UpdateGauge(metricName, value)

	case metrics.CounterType:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("cannot parse int64 from counter metric")
		}
		s.Driver.UpdateCounter(metricName, value)
	default:
		return errors.New("unknown type of metric")
	}
	return nil
}
