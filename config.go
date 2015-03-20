package coresize

import (
	"flag"
)

type Config struct {
	Port         int
	Hash         bool
	FolderName   string
	PullFrom     string // s3 or http
	PullFromUrl  string
	AwsClientKey string
	AwsSecretKey string
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) ParseFlags() {
	var (
		port         = flag.Int("port", 8080, "Port to listen on")
		hash         = flag.Bool("hash", true, "Answer to hashed filenames")
		folderName   = flag.String("folder-name", "files/", "Local folder where images to serve are located and will be pulled to")
		pullFrom     = flag.String("pull-from", "", "Either 's3' or 'http'")
		pullFromUrl  = flag.String("pull-from-url", "", "S3 location or http location")
		awsClientKey = flag.String("aws-client-key", "", "Only used when pull-from=s3")
		awsSecretKey = flag.String("aws-secret-key", "", "")
	)

	flag.Parse()

	c.Port = *port
	c.Hash = *hash
	c.FolderName = *folderName
	c.PullFrom = *pullFrom
	c.PullFromUrl = *pullFromUrl
	c.AwsClientKey = *awsClientKey
	c.AwsSecretKey = *awsSecretKey
}
