package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Env        string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		appPort = 8080
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432
	}

	return &Config{
		AppPort:    appPort,
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		Env:        os.Getenv("ENV"),
	}
}