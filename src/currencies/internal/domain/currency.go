package domain

import (
	"encoding/json"
	"log"
	"strconv"
)

type Currency struct {
	Code       string  `json:"code"`
	Codein     string  `json:"codein"`
	Name       string  `json:"name"`
	High       float64 `json:"high,omitempty"`
	Low        float64 `json:"low,omitempty"`
	VarBid     float64 `json:"varBid,omitempty"`
	PctChange  float64 `json:"pctChange,omitempty"`
	Bid        float64 `json:"bid,omitempty"`
	Ask        float64 `json:"ask,omitempty"`
	Timestamp  string  `json:"timestamp,omitempty"`
	CreateDate string  `json:"create_date,omitempty"`
}

type CurrencyData struct {
	Code      string  `json:"code"`
	Value     float64 `json:"value"`
	CreatedAt string  `json:"created_at"`
}

func NewCurrency() *Currency {
	return &Currency{}
}

func FromJSON(data string) (*Currency, error) {
	currency := &Currency{}
	err := json.Unmarshal([]byte(data), currency)
	if err != nil {
		return nil, err
	}
	return currency, nil
}

func (c *Currency) ToJSON() string {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		log.Printf("[domain.Currency] Error converting currency to json: %s", err)
		return ""
	}
	return string(data)
}

func (c *Currency) GetData() *CurrencyData {
	return &CurrencyData{
		Code:      c.Code,
		Value:     c.Bid,
		CreatedAt: c.CreateDate,
	}
}
func CurrencyFromMap(data map[string]interface{}) *Currency {
	currency := &Currency{}

	currency.Code = data["code"].(string)
	currency.Codein = data["codein"].(string)
	currency.Name = data["name"].(string)
	currency.High = GetFloat(data["high"])
	currency.Low = GetFloat(data["low"])
	currency.VarBid = GetFloat(data["varBid"])
	currency.PctChange = GetFloat(data["pctChange"])
	currency.Bid = GetFloat(data["bid"])
	currency.Ask = GetFloat(data["ask"])
	currency.Timestamp = data["timestamp"].(string)
	currency.CreateDate = data["create_date"].(string)

	return currency
}

func GetFloat(data interface{}) float64 {
	if data == nil {
		return 0
	}

	ret, err := strconv.ParseFloat(data.(string), 64)
	if err != nil {
		log.Printf("[domain.Currency] Error converting float: %s -> %s", err, data)
		return 0
	}
	return ret
}
