package domain

type Currency struct {
	Code       string  `json:"code" bson:"_id"`
	Codein     string  `json:"codein" bson:"codein"`
	Name       string  `json:"name" bson:"name"`
	High       float64 `json:"high,omitempty" bson:"high"`
	Low        float64 `json:"low,omitempty" bson:"low"`
	VarBid     float64 `json:"varBid,omitempty" bson:"varBid"`
	PctChange  float64 `json:"pctChange,omitempty" bson:"pctChange"`
	Bid        float64 `json:"bid,omitempty" bson:"bid"`
	Ask        float64 `json:"ask,omitempty" bson:"ask"`
	Timestamp  string  `json:"timestamp,omitempty" bson:"timestamp"`
	CreateDate string  `json:"create_date,omitempty" bson:"create_date"`
}
