package main

import (
	"time"

	"github.com/rombintu/yametrics/internal/agent"
	"github.com/rombintu/yametrics/internal/config"
)

func main() {
	config := config.MustLoad()
	a := agent.NewAgent(config.Agent)

	go a.RunPoll()
	go a.RunReport()
	for {
		time.Sleep(1 * time.Second)
	}
}
