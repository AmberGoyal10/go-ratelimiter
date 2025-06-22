package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvValues struct {
	REDIS_URL string
	PORT      string
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func GetEnv(key string) string {
	return goDotEnvVariable(key)
}
