package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kumparan/imaginary/config"
	log "github.com/sirupsen/logrus"
)

var (
	// S3Client :nodoc:
	S3Client *s3.S3
)

// InitializeS3Conn :nodoc:
func InitializeS3Conn() {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AWSRegion()),
		Credentials: credentials.NewStaticCredentials(config.AWSS3Key(), config.AWSS3Secret(), ""),
	})

	if err != nil {
		log.WithField("region", config.AWSRegion()).Fatalf("failed create aws session: %v", err)
	}

	log.Println("aws session created")
	S3Client = s3.New(session)
}
