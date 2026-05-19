package db

import (
	"fmt"
	"net"
)

type Config struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Database       string `yaml:"database"`
	SSLMode        string `yaml:"ssl_mode"`
	IsolationLevel string `yaml:"isolation_level"`
}

func (c *Config) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		c.User, c.Password, net.JoinHostPort(c.Host, c.Port), c.Database, c.SSLMode,
	)
}
