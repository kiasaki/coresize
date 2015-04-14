package coresize

import (
	"flag"
)

type Config struct {
	Port         int
	Verbose      bool
	AwsClientKey string
	AwsSecretKey string
	Bucket       string
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) ParseFlags() {
	var (
		port         = flag.Int("port", 8080, "Port to listen on")
		verbose      = flag.Bool("v", false, "Be more verbose")
		awsClientKey = flag.String("aws-client-key", "", "")
		awsSecretKey = flag.String("aws-secret-key", "", "")
		bucket       = flag.String("bucket", "", "S3 bucket containing images")
	)

	flag.Parse()

	c.Port = *port
	c.Verbose = *verbose
	c.AwsClientKey = *awsClientKey
	c.AwsSecretKey = *awsSecretKey
	c.Bucket = *bucket
}
