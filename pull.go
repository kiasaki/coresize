package coresize

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

func pullImages(c Config) {
	// If images folder is already present and we are not force pushing go and pull
	if _, err := os.Stat(c.FolderName); os.IsNotExist(err) || c.ForcePull {

		// If we force pull let's delete the old folder and all its contents
		if c.ForcePull {
			err = os.RemoveAll(c.FolderName)
			if err != nil {
				log.Printf("Error deleting folder: %s\n", err.Error())
				os.Exit(1)
			}
		}

		// Make sure the directory we are going to put images in exists
		err = os.MkdirAll(c.FolderName, 0777)
		if err != nil {
			log.Printf("Error creating folder: %s\n", err.Error())
			os.Exit(1)
		}

		if c.PullFrom == "s3" {
			pullFromS3(c)
		} else if c.PullFrom == "http" {
			pullFromHttp(c)
		}
	}
}

func pullFromS3(c Config) {
	log.Println("Pulling images from S3")

	awsAuth := aws.Auth{
		AccessKey: c.AwsClientKey,
		SecretKey: c.AwsSecretKey,
	}
	region := aws.USEast
	connection := s3.New(awsAuth, region)
	bucket := connection.Bucket(c.PullFromUrl)

	response, err := bucket.List("", "", "", 500)
	if err != nil {
		log.Printf("Error fetching from S3: %s\n", err.Error())
		os.Exit(1)
	}

	for _, object := range response.Contents {
		fileBytes, err := bucket.Get(object.Key)
		if err != nil {
			log.Printf("Error fetching %s from S3: %s\n", object.Key, err.Error())
			continue
		}
		err = ioutil.WriteFile(path.Join(c.FolderName, object.Key), fileBytes, 0777)
		if err != nil {
			log.Printf("Error writing %s from S3: %s\n", object.Key, err.Error())
			continue
		}
		log.Println("Downloaded: " + object.Key)
	}

	log.Println("")
}

func pullFromHttp(c Config) {
	log.Println("Pulling images from Http")

	// TODO implement http fetching of a .zip
	log.Println("Not yet implemented")
	os.Exit(1)
}
