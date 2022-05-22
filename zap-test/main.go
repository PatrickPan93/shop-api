package main

import "go.uber.org/zap"

func main() {

	logger, _ := zap.NewProduction()
	// logger, _ := zap.NewDevelopment()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

	url := "https://immoc.com"

	// 使用Sugar
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3)
	sugar.Infof("Failed to fetch URL: %s\n", url)
}
