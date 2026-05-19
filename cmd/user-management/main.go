package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/adapter/yaml"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/queue"
	"github.com/Lapakin/edu-planner/internal/user-management/repository/postgres"
	"github.com/Lapakin/edu-planner/internal/user-management/router"
	"github.com/Lapakin/edu-planner/internal/user-management/service"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
)

const defaultConfigPath = "./cmd/user-management/config.yaml"

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
	logger := logging.NewLogger(level, cfg.Logging.Formatter).WithField("service", "user-management")

	cfg.LoadDBCredentialsFromEnv()

	jwt.Init(cfg.JWT)

	logger.Infoln("Trying to connect to database...")
	db, err := pg.NewDB(cfg.DB.URL(), cfg.DB.IsolationLevel)
	if err != nil {
		logger.Fatalf("failed to connect to database. err: %v", err)
	}
	logger.Infoln("Successfully connected to database")

	awsCfg, err := cfg.SQS.LoadAWSConfig(context.TODO())
	if err != nil {
		logger.Fatalf("Failed to load AWS config for SQS producer: %v", err)
	}

	if err = queue.StartMassageProducers(cfg.SQS.QueueName, cfg.DB.URL(), db.DB, logger, awsCfg); err != nil {
		logger.Fatalf("Failed to start SQS producer: %v", err)
	}

	svc := service.NewServices(db, postgres.NewRepoManager(), logger)

	r := router.NewRouter(svc)

	logger.Infof("Starting server at %s", cfg.Addr)
	logger.Fatal(http.ListenAndServe(cfg.Addr, r))
}
