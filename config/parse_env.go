package config

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (c *Config) ParseEnv() error {
	envHTTPServHost := os.Getenv(envNameHTTPHost)
	if envHTTPServHost != "" {
		c.HTTPHost = envHTTPServHost
	}

	envHTTPServPort := os.Getenv(envNameHTTPPort)
	if envHTTPServPort != "" {
		c.HTTPPort = envHTTPServPort
	}

	c.UserCli.parseEnv()
	c.PostCli.parseEnv()
	c.RelationCli.parseEnv()
	c.FeedCli.parseEnv()
	c.MediaCli.parseEnv()

	envLogLvl := os.Getenv(envNameLogLvl)
	if envLogLvl != "" {
		logLevel, err := zerolog.ParseLevel(envLogLvl)
		if err != nil {
			return errors.Wrapf(err, "parse log lvl: %s", envLogLvl)
		}
		c.LogLvl = logLevel
	}

	return nil
}
