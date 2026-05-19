package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Lapakin/edu-planner/internal/adapter/yaml"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/ui/router"
	"github.com/Lapakin/edu-planner/internal/ui/service"
)

const defaultConfigPath = "./cmd/ui/config.yaml"

func main() {
	var configPath string
	flag.StringVar(
		&configPath,
		"config",
		defaultConfigPath,
		"path to config file",
	)
	flag.Parse()

	var cfg config
	if err := yaml.UnmarshalYAML(configPath, &cfg); err != nil {
		log.Fatalf("failed to unmarshal configuration. err: %v", err)
	}

	level, err := logging.ParseLevel(cfg.Logging.Level)
	if err != nil {
		log.Fatalf("failed to parse logging level. err: %v", err)
	}
	logger := logging.NewLogger(level, cfg.Logging.Formatter).WithField("service", "ui")

	svc := service.NewServices(logger, cfg.UI.UserManagement, cfg.UI.Syllabus)

	r := router.NewRouter(svc)

	logger.Infof("starting server on %s", cfg.Addr)
	logger.Fatal(http.ListenAndServe(cfg.Addr, r))
}
