package main

import (
	"os"

	"github.com/Lapakin/edu-planner/internal/adapter/broker"
	"github.com/Lapakin/edu-planner/internal/adapter/db"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/logging"
)

type config struct {
	Addr    string          `yaml:"addr"`
	DB      *db.Config      `yaml:"db"`
	Logging *logging.Config `yaml:"logging"`
	SQS     *broker.Config  `yaml:"sqs"`
	JWT     *jwt.Config     `yaml:"jwt"`
}

func (c *config) LoadDBCredentialsFromEnv() {
	if envUser := os.Getenv("DB_USER"); envUser != "" {
		c.DB.User = envUser
	}
	if envPass := os.Getenv("DB_PASSWORD"); envPass != "" {
		c.DB.Password = envPass
	}
}
