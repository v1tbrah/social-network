package config

import (
	"github.com/rs/zerolog"
)

type Config struct {
	HTTPServHost string
	HTTPServPort string

	UserServiceClient     UserServiceClient
	PostServiceClient     PostServiceClient
	RelationServiceClient RelationServiceClient
	FeedServiceClient     FeedServiceClient

	LogLvl zerolog.Level
}

func NewDefaultConfig() Config {
	return Config{
		HTTPServHost: defaultHTTPServHost,
		HTTPServPort: defaultHTTPServPort,

		UserServiceClient:     newDefaultUserServiceClientConfig(),
		PostServiceClient:     newDefaultPostServiceClientConfig(),
		RelationServiceClient: newDefaultRelationServiceClientConfig(),
		FeedServiceClient:     newDefaultFeedServiceClientConfig(),

		LogLvl: defaultLogLvl,
	}
}
