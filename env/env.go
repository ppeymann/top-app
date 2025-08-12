package env

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key, def string) string {
	err := godotenv.Load()
	if err != nil {
		return def
	}

	val := os.Getenv(string(key))
	if len(val) == 0 {
		return def
	}

	return val
}

func IsProduction() bool {
	mode := GetEnv("GIN_MODE", "debug")
	return mode == "release"
}
