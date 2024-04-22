package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost           string
	Port                 string
	DBUser               string
	DBPassword           string
	DBAddress            string
	DBName               string
	Net                  string
	AllowNativePasswords bool
	ParseTime            bool
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:           getEnv("PUBLIC_HOST", "localhost"),
		Port:                 getEnv("PORT", "8080"),
		DBUser:               getEnv("DB_USER", "root"),
		DBPassword:           getEnv("DB_PASSWORD", "root"),
		DBAddress:            fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "3307")),
		DBName:               getEnv("DB_NAME", "scootin_aboot"),
		Net:                  getEnv("DB_NET", "tcp"),
		AllowNativePasswords: getEnv("DB_ALLOW_NATIVE_PASSWORDS", "true") == "true",
		ParseTime:            getEnv("DB_PARSE_TIME", "true") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
