package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderProducts []*Product         `json:"orders_products,omitempty" bson:"orders_products,omitempty"`
	TotalPrice    float64            `json:"total_price,omitempty" bson:"total_price,omitempty"`
	RealPrice     float64            `json:"real_price,omitempty" bson:"real_price,omitempty"` //Satış fiyatı
	PaymentMethod uint               `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	SalesType     uint               `json:"sales_type,,omitempty" bson:"sales_type,omitempty"` //Barcode numarası
}
