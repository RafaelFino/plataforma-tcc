package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Debug              bool   `json:"debug,omitempty"`
	ServerAddress      string `json:"server_address"`
	ServerPort         int    `json:"server_port"`
	LogPath            string `json:"log_path"`
	DBURL              string `json:"db_url"`
	DBName             string `json:"db_name"`
	ProductsURL        string `json:"products_url"`
	CurrenciesURL      string `json:"currencies_url"`
	ClientsURL         string `json:"clients_url"`
	CurrenciesInterval int    `json:"currencies_interval"`
}

func NewConfig() *Config {
	return &Config{}
}

func ConfigFromJSON(data string) (*Config, error) {
	config := &Config{}
	err := json.Unmarshal([]byte(data), config)
	if err != nil {
		return nil, err
	}

	log.Printf("[config] loaded config: %+v", config)

	return config, nil
}

func (c *Config) ToJSON() string {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return ""
	}
	return string(data)
}

func ConfigClientFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ConfigFromJSON(string(data))
}
