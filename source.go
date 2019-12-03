package main

import (
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/service/s3"
)

type ImageSourceType string
type ImageSourceFactoryFunction func(*SourceConfig) ImageSource

type SourceConfig struct {
	AuthForwarding bool
	Authorization  string
	MountPath      string
	Type           ImageSourceType
	ForwardHeaders []string
	AllowedOrigins []*url.URL
	MaxAllowedSize int
	S3Client       *s3.S3
}

var imageSourceMap = make(map[ImageSourceType]ImageSource)
var imageSourceFactoryMap = make(map[ImageSourceType]ImageSourceFactoryFunction)

type ImageSource interface {
	Matches(*http.Request) bool
	GetImage(*http.Request) ([]byte, error)
}

func RegisterSource(sourceType ImageSourceType, factory ImageSourceFactoryFunction) {
	imageSourceFactoryMap[sourceType] = factory
}

func LoadSources(o ServerOptions) {
	for name, factory := range imageSourceFactoryMap {
		imageSourceMap[name] = factory(&SourceConfig{
			Type:           name,
			MountPath:      o.Mount,
			AuthForwarding: o.AuthForwarding,
			Authorization:  o.Authorization,
			AllowedOrigins: o.AllowedOrigins,
			MaxAllowedSize: o.MaxAllowedSize,
			ForwardHeaders: o.ForwardHeaders,
			S3Client:       o.S3Client,
		})
	}
}

func MatchSource(req *http.Request) ImageSource {
	for _, source := range imageSourceMap {
		if source.Matches(req) {
			return source
		}
	}
	return nil
}
