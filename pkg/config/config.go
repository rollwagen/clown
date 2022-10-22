package config

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

type Config struct {
	AuthToken string
	Host      string
}

func New() *Config {
	homeDir, _ := os.UserHomeDir()

	_ = godotenv.Load(homeDir + "/.clown")

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	gitlabHost := os.Getenv("GITLAB_HOST")

	log.Debug().Msgf("Using GITLAB_TOKEN: %s GITLAB_HOST: %s", gitlabToken, gitlabHost)

	return &Config{
		AuthToken: gitlabToken,
		Host:      gitlabHost,
	}
}
