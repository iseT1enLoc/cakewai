package appconfig

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		//fmt.Errorf("Could not load .env %v", err)
		return nil, err
	}
	acc_time := os.Getenv("ACCESS_TOK_EXP")
	re_time := os.Getenv("REFRESH_TOK_EXP")
	// Convert string to integer
	iacc_time, err1 := strconv.Atoi(acc_time)
	ire_time, err2 := strconv.Atoi(re_time)
	if err1 != nil || err2 != nil {
		// Handle the error if conversion fails
		fmt.Println("Error converting string to integer:", err1)
		return nil, err
	}
	return &Env{
		DB_USER:         os.Getenv("DB_USER"),
		DB_PASSWORD:     os.Getenv("DB_PASSWORD"),
		DB_HOST:         os.Getenv("DB_HOST"),
		DB_NAME:         os.Getenv("DB_NAME"),
		DB_PORT:         os.Getenv("DB_PORT"),
		SECRET_KEY:      os.Getenv("SECRET_KEY"),
		ACCESS_SECRET:   os.Getenv("ACCESS_KEY"),
		REFRESH_SECRET:  os.Getenv("REFRESH_KEY"),
		TIMEOUT:         os.Getenv("TIMEOUT"),
		ACCESS_TOK_EXP:  iacc_time,
		REFRESH_TOK_EXP: ire_time,
	}, nil
}
