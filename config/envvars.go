package config

import (
	"log"
	"os"
)

func RequireEnvVar(key string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		log.Fatalf("%s env var required!", key)
	}

	return value
}
