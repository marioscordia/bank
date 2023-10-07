package config

import (
	"encoding/json"
	"os"
)

type DB struct {
	Driver string `json:"driver"`
	Name string `json:"name"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	Server string `json:"server"`
	Secret string `json:"secret"`
	DB DB `json:"db"`
}

func NewConfig() (*Config, error) {
		file, err := os.Open("./config/config.json")
		if err != nil {
			return nil, err
		}
		defer file.Close()
	
		var config *Config
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			return nil, err
		}

		if err := os.Setenv("JWT_SECRET", config.Secret); err != nil {
			return nil, err
		}
		return config, nil
}