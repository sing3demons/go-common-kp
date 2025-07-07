package main

import (
	"os"

	config "github.com/sing3demons/go-common-kp/kp/configs"
	"github.com/sing3demons/go-common-kp/kp/pkg/kp"
)

func main() {
	conf := config.NewConfig()
	if os.Getenv("ENV") == "docker" {
		conf.LoadEnv("configs/.docker.env")
	} else {
		conf.LoadEnv("configs")
	}

	app := kp.NewApplication(conf)
	// app.StartKafka()

	app.Get("/healthz", func(ctx *kp.Context) error {
		ctx.Info("Health check endpoint hit")
		return ctx.JSON(200, "OK")
	})

	app.Start()
}
