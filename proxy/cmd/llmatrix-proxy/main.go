package main

import (
	"log"
	"os"

	"github.com/icosmos-space/llmatrix/proxy/internal/bootstrap"
	"github.com/icosmos-space/llmatrix/proxy/internal/config"
)

func main() {
	cfgPath := os.Getenv("LLMATRIX_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/config.example.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}
