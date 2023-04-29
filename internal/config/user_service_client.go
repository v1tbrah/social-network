package config

import "os"

const (
	defaultUserServiceClientHost = "127.0.0.1"
	defaultUserServiceClientPort = "6060"
)

const (
	envNameUserServiceClientHost = "USER_SERVICE_CLIENT_HOST"
	envNameUserServiceClientPort = "USER_SERVICE_CLIENT_PORT"
)

type UserServiceClient struct {
	ServHost string
	ServPort string
}

func newDefaultUserServiceClientConfig() UserServiceClient {
	return UserServiceClient{
		ServHost: defaultUserServiceClientHost,
		ServPort: defaultUserServiceClientPort,
	}
}

func (u *UserServiceClient) parseEnv() {
	envServHost := os.Getenv(envNameUserServiceClientHost)
	if envServHost != "" {
		u.ServHost = envServHost
	}

	envServPort := os.Getenv(envNameUserServiceClientPort)
	if envServPort != "" {
		u.ServPort = envServPort
	}
}
