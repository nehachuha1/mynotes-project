package main

import "go.uber.org/zap"

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugaredLogger := logger.Sugar()
	sugaredLogger.Infow("started server", "type", "START", "addr", "test")
}
