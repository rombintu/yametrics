package agent

import (
	"testing"

	"github.com/rombintu/yametrics/internal/config"
)

func TestAgentSendData(t *testing.T) {

	config := config.AgentConfig{
		ServerUrl:      "http://localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
		Mode:           "debug",
	}
	agent := NewAgent(config)
	agent.sendDataOnServer(
		"counter",
		"c1",
		"100",
	)

}
