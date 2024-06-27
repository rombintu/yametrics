package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rombintu/yametrics/internal/config"
	"github.com/rombintu/yametrics/internal/server"
)

func main() {

	config := config.MustLoad()

	server := server.NewServer(config)
	server.Log.Info("starting application", slog.String("env", config.Environment))
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go server.Start()

	sign := <-stop
	server.Shutdown()
	server.Log.Info("server stopped", slog.String("signal", sign.String()))
}
