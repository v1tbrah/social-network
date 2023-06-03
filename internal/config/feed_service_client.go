package config

import "os"

const (
	defaultFeedServiceClientHost = "127.0.0.1"
	defaultFeedServiceClientPort = "3030"
)

const (
	envNameFeedServiceClientHost = "FEED_SERVICE_CLIENT_HOST"
	envNameFeedServiceClientPort = "FEED_SERVICE_CLIENT_PORT"
)

type FeedServiceClient struct {
	Host string
	Port string
}

func newDefaultFeedServiceClientConfig() FeedServiceClient {
	return FeedServiceClient{
		Host: defaultFeedServiceClientHost,
		Port: defaultFeedServiceClientPort,
	}
}

func (u *FeedServiceClient) parseEnv() {
	envHost := os.Getenv(envNameFeedServiceClientHost)
	if envHost != "" {
		u.Host = envHost
	}

	envPort := os.Getenv(envNameFeedServiceClientPort)
	if envPort != "" {
		u.Port = envPort
	}
}
