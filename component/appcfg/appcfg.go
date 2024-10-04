package appconfig

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		//fmt.Errorf("Could not load .env %v", err)
		return nil, err
	}
	return &Env{
		DB_USER:        os.Getenv("DB_USER"),
		DB_PASSWORD:    os.Getenv("DB_PASSWORD"),
		DB_HOST:        os.Getenv("DB_HOST"),
		DB_NAME:        os.Getenv("DB_NAME"),
		DB_PORT:        os.Getenv("DB_PORT"),
		SECRET_KEY:     os.Getenv("SECRET_KEY"),
		ACCESS_SECRET:  os.Getenv("ACCESS_KEY"),
		REFRESH_SECRET: os.Getenv("REFRESH_KEY"),
		TIMEOUT:        os.Getenv("TIMEOUT"),
	}, nil
}
