package configs

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"os"
)

const configPath = "./configs/config.json"

type Config struct {
	Server struct {
		Host string `validate:"required"`
	}
	Postgres struct {
		Host     string `json:"Host"`
		Port     int    `json:"Port"`
		User     string `json:"User"`
		Password string `json:"Password"`
		DBName   string `json:"DBName"`
	}
	SecretKey string `json:"SecretKey"`
}

func LoadConfig() (c *Config, err error) {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(jsonFile).Decode(&c)
	if err != nil {
		return nil, err
	}

	err = validator.New().Struct(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
