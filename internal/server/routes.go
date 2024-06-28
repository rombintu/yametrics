package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/rombintu/yametrics/internal/metrics"
)

func rootHandler(c echo.Context) error {
	return c.String(http.StatusOK, pongMessage)
}

func (s *Server) updateMetrics(c echo.Context) error {

	mtype := c.Param("mtype")
	mname := c.Param("mname")
	mvalue := c.Param("mvalue")

	// Check mtype is counter or gauge
	if mtype != metrics.CounterType && mtype != metrics.GaugeType {
		return c.String(http.StatusNotFound, "invalid metric type")
	}

	// Check mvalue is Float64 or int64
	if _, err := strconv.ParseFloat(mvalue, 64); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if _, err := strconv.ParseInt(mvalue, 10, 64); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Write metric to storage
	if err := s.storage.WriteMetric(mtype, mname, mvalue); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "updated")

}

// GET http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func (s *Server) getMetricsByNames(c echo.Context) error {
	mtype := c.Param("mtype")
	mname := c.Param("mname")
	value := s.storage.Driver.GetMetricByName(mtype, mname)
	if value == "" {
		return c.String(http.StatusNotFound, "metric not found")
	}
	return c.String(http.StatusOK, value)
}
