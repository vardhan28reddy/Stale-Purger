package main

import (
	cfg "Stale-purger/pkg/config"
	"Stale-purger/pkg/controller"
	"Stale-purger/pkg/utils"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize config
	config, err := cfg.InitializeConfig()
	if err != nil {
		panic(err)
	}
	// Initialize logger
	logger := cfg.InitializeLogger(config.LogLevel, config.LogJsonFormat)

	logger.Info("Initializing Postgres DB connection...")
	postgresDB, err := config.InitializeDB(logger)
	utils.FatalFunc("Couldn't connect Postgres DB", err, logger)
	logger.Info("Initialized DB connection pool")

	logger.Info("Initializing K8s client...")
	kubeClient, err := config.InitializeKubeClient(logger)
	utils.FatalFunc("Couldn't create K8s client", err, logger)
	logger.Info("Initialized K8s client")

	ctx, cancel := context.WithCancel(context.Background())

	c := controller.NewController(*config, logger)
	c.AddComponent(controller.NewPurgerComponent(postgresDB, logger, kubeClient, *config))
	c.Start(ctx)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	logger.Info("Received shutdown signal")
	cancel()
	logger.Info("Closing DB connection")
	if err := postgresDB.Close(); err != nil {
		logger.Error("Couldn't close DB gracefully")
	}
}
