package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt     time.Time          `json:"order_created_date,omitempty" bson:"order_created_date,omitempty"`
	OrderProducts []Product          `json:"order_products,omitempty" bson:"order_products,omitempty"`
	TotalPrice    int                `json:"total_price,omitempty" bson:"total_price,omitempty"`
	RealPrice     float64            `json:"real_price,omitempty" bson:"real_price,omitempty"` //Satış fiyatı
	PaymentMethod string             `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	SalesType     string             `json:"sales_type,omitempty" bson:"sales_type,omitempty"` //Barcode numarası
}
