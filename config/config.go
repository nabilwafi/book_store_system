package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     int
	GRPCPort    int
	JWTSecret   string
	JWTExpiry   int
	MidtransKey string
	DBConfig    DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	appPort, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))
	grpcPort, _ := strconv.Atoi(getEnv("GRPC_PORT", "50051"))
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY", "24"))

	return &Config{
		AppPort:     appPort,
		GRPCPort:    grpcPort,
		JWTSecret:   getEnv("JWT_SECRET", "hello_123"),
		JWTExpiry:   jwtExpiry,
		MidtransKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		DBConfig: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "bookstore"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
