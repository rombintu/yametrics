package main

import (
	"github.com/rombintu/yametrics/internal/agent"
	"github.com/rombintu/yametrics/internal/config"
)

func main() {
	config := config.MustLoad()
	a := agent.NewAgent(config.Agent)
	a.Run()
}
