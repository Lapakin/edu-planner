package broker

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type Config struct {
	QueueName string `yaml:"queue_name"`
	Endpoint  string `yaml:"endpoint"`
	Region    string `yaml:"region"`
}

func (c *Config) LoadAWSConfig(ctx context.Context) (aws.Config, error) {
	if c == nil {
		return aws.Config{}, errors.New("SQS Config is nil, it was not loaded from yaml")
	}
	opts := []func(*config.LoadOptions) error{
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	}
	if c.Region != "" {
		opts = append(opts, config.WithRegion(c.Region))
	}
	if c.Endpoint != "" {
		opts = append(opts, config.WithBaseEndpoint(c.Endpoint))
	}
	return config.LoadDefaultConfig(ctx, opts...)
}
