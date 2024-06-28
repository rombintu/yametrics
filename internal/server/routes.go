package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/rombintu/yametrics/internal/metrics"
)

func (s *Server) updateMetrics(c echo.Context) error {

	mtype := c.Param("mtype")
	mname := c.Param("mname")
	mvalue := c.Param("mvalue")

	// Check mvalue is Float64 or int64
	switch mtype {
	case metrics.CounterType:
		if _, err := strconv.ParseInt(mvalue, 10, 64); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	case metrics.GaugeType:
		if _, err := strconv.ParseFloat(mvalue, 64); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	default:
		return c.String(http.StatusBadRequest, "invalid metric type")
	}

	// Write metric to storage
	if err := s.storage.WriteMetric(mtype, mname, mvalue); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "updated")

}

// GET http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s *Server) getMetricsByNames(c echo.Context) error {
	mtype := c.Param("mtype")
	mname := c.Param("mname")
	value := s.storage.Driver.GetMetricByName(mtype, mname)
	if value == "" {
		return c.String(http.StatusNotFound, "not found")
	}
	return c.String(http.StatusOK, value)
}
