package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
    // "github.com/aws/aws-sdk-go/service/s3"
)


func ConfigAWSSession() (*session.Session, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	AccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	SecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, ""),
	})

	if err != nil {
		log.Fatal(err)
	}

	return sess, nil
}
