package main

import (
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/ui"
)

type config struct {
	Addr    string          `yaml:"addr"`
	Logging *logging.Config `yaml:"logging"`
	UI      ui.Config       `yaml:"ui"`
}
