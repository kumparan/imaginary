package main

import (
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kumparan/imaginary/config"
	log "github.com/sirupsen/logrus"
)

const ImageSourceTypeS3 ImageSourceType = "s3"
const S3QueryKey = "s3"

type S3ImageSource struct {
	Config *SourceConfig
}

func NewS3ImageSource(config *SourceConfig) ImageSource {
	return &S3ImageSource{
		Config: config,
	}
}

func (s *S3ImageSource) Matches(r *http.Request) bool {
	return r.Method == http.MethodGet && r.URL.Query().Get(S3QueryKey) != ""
}

func (s *S3ImageSource) GetImage(req *http.Request) ([]byte, error) {
	awsS3Key := parseS3Path(req)
	return s.fetchImage(awsS3Key, req)
}

func (s *S3ImageSource) fetchImage(awsS3Key string, ireq *http.Request) (b []byte, err error) {
	// Check remote image size by fetching HTTP Headers
	result, err := s.Config.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.AWSS3Bucket()),
		Key:    aws.String(awsS3Key),
	})

	if err != nil {
		log.WithFields(log.Fields{
			"awsS3Key": awsS3Key}).
			Info(err)
		return nil, err
	}

	b, err = ioutil.ReadAll(result.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"awsS3Key": awsS3Key}).
			Error(err)
		return nil, err
	}

	return
}

func parseS3Path(request *http.Request) (awsS3Key string) {
	return request.URL.Query().Get(S3QueryKey)
}

func init() {
	RegisterSource(ImageSourceTypeS3, NewS3ImageSource)
}
