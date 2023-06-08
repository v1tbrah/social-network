package config

import "os"

const (
	defaultRelationServiceClientHost = "127.0.0.1"
	defaultRelationServiceClientPort = "4040"
)

const (
	envNameRelationServiceClientHost = "RELATION_SERVICE_CLIENT_HOST"
	envNameRelationServiceClientPort = "RELATION_SERVICE_CLIENT_PORT"
)

type RelationCli struct {
	Host string
	Port string
}

func newDefaultRelationServiceClientConfig() RelationCli {
	return RelationCli{
		Host: defaultRelationServiceClientHost,
		Port: defaultRelationServiceClientPort,
	}
}

func (u *RelationCli) parseEnv() {
	envHost := os.Getenv(envNameRelationServiceClientHost)
	if envHost != "" {
		u.Host = envHost
	}

	envPort := os.Getenv(envNameRelationServiceClientPort)
	if envPort != "" {
		u.Port = envPort
	}
}
