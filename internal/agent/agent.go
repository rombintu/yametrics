package agent

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

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
	pollInterval   int64
	reportInterval int64
	log            *slog.Logger
	data           map[string]interface{}
	pollCount      int
	metrics        map[string]string
}

func NewAgent(config config.AgentConfig) *Agent {
	data := make(map[string]interface{})
	data["counter"] = make(storage.CounterTable)
	data["gauge"] = make(storage.GaugeTable)
	return &Agent{
		serverUrl:      fixServerUrl(config.ServerUrl),
		pollInterval:   config.PollInterval,
		reportInterval: config.ReportInterval,
		log:            lib.SetupLogger(config.Mode),
		data:           data,
		metrics:        make(map[string]string),
	}
}

func fixServerUrl(url string) string {
	if strings.HasPrefix(url, "http://") {
		return url
	} else {
		return fmt.Sprintf("http://%s", url)
	}
}

func (a *Agent) incPollCount() {
	a.pollCount++
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

func (a *Agent) sendDataOnServer(metricType, metricName string, value string) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", a.serverUrl, metricType, metricName, value)
	a.log.Debug(url)
	if err := a.postRequest(url); err != nil {
		return err
	}
	return nil
}

func (a *Agent) RunPoll() {
	a.log.Info("start poll", slog.Any("pollInterval", a.pollInterval))
	for {
		a.loadMetrics()
		time.Sleep(time.Duration(a.pollInterval) * time.Second)
	}
}

func (a *Agent) RunReport() {
	a.log.Info("start report", slog.Any("reportInterval", a.reportInterval))
	for {
		a.log.Debug("New report", slog.Int("pollCount", a.pollCount))
		for metricName, value := range a.metrics {
			a.sendDataOnServer(gaugeMetricType, metricName, value)
		}
		a.sendDataOnServer(counterMetricType, "pollCount", strconv.Itoa(a.pollCount))
		time.Sleep(time.Duration(a.reportInterval) * time.Second)
	}
}

func (a *Agent) loadMetrics() {
	a.log.Debug("load metrics")
	var metrics runtime.MemStats
	runtime.ReadMemStats(&metrics)

	var metricsInterface map[string]interface{}
	inrec, err := json.Marshal(metrics)
	if err != nil {
		a.log.Error(err.Error())
		return
	}
	json.Unmarshal(inrec, &metricsInterface)
	for name, value := range metricsInterface {
		switch v := value.(type) {
		case int:
			a.metrics[name] = strconv.Itoa(v)
		case float64:
			a.metrics[name] = strconv.FormatFloat(v, 'f', -1, 64)
		case uint64:
			a.metrics[name] = strconv.FormatUint(v, 10)
		}

	}

	// get random float64 value
	randomValue := rand.Float64()
	a.metrics["randomValue"] = strconv.FormatFloat(randomValue, 'f', -1, 64)
	a.log.Debug("load metrics is complete", slog.Int("pollCount", a.pollCount))
	a.incPollCount()
}
