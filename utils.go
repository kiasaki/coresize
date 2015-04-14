package coresize

import (
	"os"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

func ensureFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func s3BucketFromConfig(c Config) *s3.Bucket {
	awsAuth := aws.Auth{
		AccessKey: c.AwsClientKey,
		SecretKey: c.AwsSecretKey,
	}
	region := aws.USEast
	connection := s3.New(awsAuth, region)
	return connection.Bucket(c.Bucket)
}
