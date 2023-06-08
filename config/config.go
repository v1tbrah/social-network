package config

import (
	"github.com/rs/zerolog"
)

type Config struct {
	HTTPHost string
	HTTPPort string

	UserCli     UserCli
	PostCli     PostCli
	RelationCli RelationCli
	FeedCli     FeedCli

	LogLvl zerolog.Level
}

func NewDefaultConfig() Config {
	return Config{
		HTTPHost: defaultHTTPHost,
		HTTPPort: defaultHTTPPort,

		UserCli:     newDefaultUserServiceClientConfig(),
		PostCli:     newDefaultPostServiceClientConfig(),
		RelationCli: newDefaultRelationServiceClientConfig(),
		FeedCli:     newDefaultFeedServiceClientConfig(),

		LogLvl: defaultLogLvl,
	}
}
