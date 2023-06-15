package config

import "os"

const (
	defaultMediaServiceClientHost = "127.0.0.1"
	defaultMediaServiceClientPort = "2020"
)

const (
	envNameMediaServiceClientHost = "MEDIA_SERVICE_CLIENT_HOST"
	envNameMediaServiceClientPort = "MEDIA_SERVICE_CLIENT_PORT"
)

type MediaCli struct {
	Host string
	Port string
}

func newDefaultMediaServiceClientConfig() MediaCli {
	return MediaCli{
		Host: defaultMediaServiceClientHost,
		Port: defaultMediaServiceClientPort,
	}
}

func (u *MediaCli) parseEnv() {
	envHost := os.Getenv(envNameMediaServiceClientHost)
	if envHost != "" {
		u.Host = envHost
	}

	envPort := os.Getenv(envNameMediaServiceClientPort)
	if envPort != "" {
		u.Port = envPort
	}
}
