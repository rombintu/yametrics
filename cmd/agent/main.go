package main

import (
	"time"

	"github.com/rombintu/yametrics/internal/agent"
	"github.com/rombintu/yametrics/internal/config"
)

func main() {
	config := config.LoadAgentConfigFromFlags()
	a := agent.NewAgent(config)

	go a.RunPoll()
	go a.RunReport()
	for {
		time.Sleep(1 * time.Second)
	}
}
