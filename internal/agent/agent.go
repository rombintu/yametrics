package agent

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/rombintu/yametrics/internal/config"
	"github.com/rombintu/yametrics/internal/storage"
	"github.com/rombintu/yametrics/lib"
)

const (
	counterMetricType = "counter"
	gaugeMetricType   = "gauge"
)

type Agent struct {
	serverUrl      string
	pollInterval   int
	reportInterval int
	log            *slog.Logger
	metrics        runtime.MemStats
	data           map[string]interface{}
}

func NewAgent(config config.AgentConfig) *Agent {
	data := make(map[string]interface{})
	data["counter"] = make(storage.CounterTable)
	data["gauge"] = make(storage.GaugeTable)
	return &Agent{
		serverUrl:      config.ServerUrl,
		pollInterval:   config.PollInterval,
		reportInterval: config.ReportInterval,
		log:            lib.SetupLogger(config.Mode),
		metrics:        runtime.MemStats{},
		data:           data,
	}
}

func (a *Agent) Run() {
	// for {
	// 	a.poll()
	// 	a.report()
	// }
	a.log.Error("implement me")
}

func (a *Agent) postRequest(url string) error {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	a.log.Debug(url, slog.Int("status", resp.StatusCode))
	defer resp.Body.Close()
	return nil
}

func (a *Agent) sendDataOnServer(metricType, metricName string, value interface{}) error {
	url := fmt.Sprintf("%s/%s/%s/%d", a.serverUrl, metricType, metricName, value)
	if err := a.postRequest(url); err != nil {
		return err
	}
	return nil
}

func (a *Agent) poll() {

}

func (a *Agent) report() {

}

func (a *Agent) loadMetrics() {
	a.log.Debug("load metrics")
	runtime.ReadMemStats(&a.metrics)

	a.log.Debug("load metrics is complete")
}
