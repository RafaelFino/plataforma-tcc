package domain

import "encoding/json"

type Currency struct {
	Code       string  `json:"code"`
	Codein     string  `json:"codein"`
	Name       string  `json:"name"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	VarBid     float64 `json:"varBid"`
	PctChange  float64 `json:"pctChange"`
	Bid        float64 `json:"bid"`
	Ask        float64 `json:"ask"`
	Timestamp  string  `json:"timestamp"`
	CreateDate string  `json:"create_date"`
}

func NewCurrency() *Currency {
	return &Currency{}
}

func FromJSON(data string) (*Currency, error) {
	currency := &Currency{}
	err := json.Unmarshal(data, currency)
	if err != nil {
		return nil, err
	}
	return currency, nil
}

func FromJSONList(data string) (map[string]*Currency, error) {
	currencies := make(map[string]*Currency)
	err := json.Unmarshal([]byte(data), &currencies)
	if err != nil {
		return nil, err
	}
	return currencies, nil
}

func (c *Currency) ToJSON() (string, error) {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
