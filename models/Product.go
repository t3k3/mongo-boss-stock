package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Detail       string             `json:"detail,omitempty" bson:"detail,omitempty"`
	Price        float64            `json:"price,omitempty" bson:"price,omitempty"` //Satış fiyatı
	Quantity     int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Barcode      uint               `json:"barcode,,omitempty" bson:"barcode,omitempty"`            //Barcode numarası
	StoreID      uint               `json:"store_id,omitempty" bson:"store_id,omitempty"`           //StoreID
	CategoryName string             `json:"category_name,omitempty" bson:"category_name,omitempty"` //Ürünün kaegorisi
	Entry_Price  float64            `json:"entry_price,omitempty" bson:"entry_price,omitempty"`     //Alış fiyatı
	Tax          float64            `json:"kdv,omitempty" bson:"kdv,omitempty"`                     //Vergi oranı
}
