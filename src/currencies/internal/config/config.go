package Config

import "encoding/json"

type Config struct {
	Debug         bool   `json:"debug,omitempty"`
	CurrencyURL   string `json:"currency_url"`
	ServerAddress string `json:"server_address"`
	ServerPort    int    `json:"server_port"`
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
	return config, nil
}

func (c *Config) ToJSON() (string, error) {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
