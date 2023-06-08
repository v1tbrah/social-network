package config

import "github.com/rs/zerolog"

const (
	defaultHTTPHost = "127.0.0.1"
	defaultHTTPPort = "8080"
	defaultLogLvl   = zerolog.InfoLevel
)

const (
	envNameHTTPHost = "HTTP_HOST"
	envNameHTTPPort = "HTTP_PORT"
	envNameLogLvl   = "LOG_LVL"
)
