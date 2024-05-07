package domain

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
