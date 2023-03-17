package model

type Brand struct {
	ID           int64   `json:"id" bson:"id"`
	Name         string  `json:"brand_name" bson:"brand_name"`
	AveragePrice float64 `json:"average_price" bson:"average_price"`
	Models       []Model `json:"models" bson:"models"`
}

type Model struct {
	ID           int64   `json:"id" bson:"id"`
	Name         string  `json:"name" bson:"name"`
	AveragePrice float64 `json:"average_price" bson:"average_price"`
	BrandName    string  `json:"brand_name" bson:"brand_name"`
}
