package config

import "github.com/rs/zerolog"

const (
	defaultHTTPServHost = "127.0.0.1"
	defaultHTTPServPort = "8080"
	defaultLogLvl       = zerolog.InfoLevel
)

const (
	envNameHTTPServHost = "HTTP_HOST"
	envNameHTTPServPort = "HTTP_PORT"
	envNameLogLvl       = "LOG_LVL"
)
