package configuration

import (
	"log"
	"os"
	"strconv"

	"github.com/casmeyu/micro-user/structs"
	"github.com/joho/godotenv"
)

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func LoadConfig(cfg *structs.Config) error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("[config] (loadConfig) - Error while loading .env file: ", err.Error())
		return err
	}

	cfg.App = structs.AppConfig{
		Name: os.Getenv("APP_NAME"),
		Ip:   os.Getenv("APP_IP"),
	}

	cfg.Db = structs.DbConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		Ip:       os.Getenv("DB_IP"),
	}

	cfg.Jwt = structs.JwtConfig{
		Secret:     os.Getenv("JWT_SECRET"),
		Expiration: getEnvAsInt("JWT_EXPIRATION", 200000),
	}

	return nil
}
