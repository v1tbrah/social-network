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

type PostServiceClient struct {
	ServHost string
	ServPort string
}

func newDefaultPostServiceClientConfig() PostServiceClient {
	return PostServiceClient{
		ServHost: defaultPostServiceClientHost,
		ServPort: defaultPostServiceClientPort,
	}
}

func (u *PostServiceClient) parseEnv() {
	envServHost := os.Getenv(envNamePostServiceClientHost)
	if envServHost != "" {
		u.ServHost = envServHost
	}

	envServPort := os.Getenv(envNamePostServiceClientPort)
	if envServPort != "" {
		u.ServPort = envServPort
	}
}
