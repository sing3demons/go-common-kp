package main

import (
	"os"

	config "github.com/sing3demons/go-common-kp/kp/configs"
	"github.com/sing3demons/go-common-kp/kp/pkg/kp"
	"github.com/sing3demons/go-common-kp/kp/pkg/logger"
)

func main() {
	conf := config.NewConfig()
	if os.Getenv("ENV") == "docker" {
		conf.LoadEnv("configs/.docker.env")
	} else {
		conf.LoadEnv("configs")
	}

	logApp := logger.NewLogger(conf.Log.App)
	defer logApp.Sync()

	logDetail := logger.NewLogger(conf.Log.Detail)
	defer logDetail.Sync()
	logSummary := logger.NewLogger(conf.Log.Summary)
	defer logSummary.Sync()

	app := kp.NewApplication(conf, logApp)
	app.LogDetail(logDetail)
	app.LogSummary(logSummary)

	app.Get("/healthz", func(ctx *kp.Context) error {
		return ctx.JSON(200, "OK")
	})

	app.Start()
}
