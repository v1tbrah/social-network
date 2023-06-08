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

type UserCli struct {
	Host string
	Port string
}

func newDefaultUserServiceClientConfig() UserCli {
	return UserCli{
		Host: defaultUserServiceClientHost,
		Port: defaultUserServiceClientPort,
	}
}

func (u *UserCli) parseEnv() {
	envHost := os.Getenv(envNameUserServiceClientHost)
	if envHost != "" {
		u.Host = envHost
	}

	envPort := os.Getenv(envNameUserServiceClientPort)
	if envPort != "" {
		u.Port = envPort
	}
}
