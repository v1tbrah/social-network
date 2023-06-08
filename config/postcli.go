package config

import "os"

const (
	defaultPostServiceClientHost = "127.0.0.1"
	defaultPostServiceClientPort = "5050"
)

const (
	envNamePostServiceClientHost = "POST_SERVICE_CLIENT_HOST"
	envNamePostServiceClientPort = "POST_SERVICE_CLIENT_PORT"
)

type PostCli struct {
	Host string
	Port string
}

func newDefaultPostServiceClientConfig() PostCli {
	return PostCli{
		Host: defaultPostServiceClientHost,
		Port: defaultPostServiceClientPort,
	}
}

func (u *PostCli) parseEnv() {
	envHost := os.Getenv(envNamePostServiceClientHost)
	if envHost != "" {
		u.Host = envHost
	}

	envPort := os.Getenv(envNamePostServiceClientPort)
	if envPort != "" {
		u.Port = envPort
	}
}
