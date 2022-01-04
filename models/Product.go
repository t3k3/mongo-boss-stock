package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	Detail        string             `json:"detail,omitempty" bson:"detail,omitempty"`
	Price         float64            `json:"price,omitempty" bson:"price,omitempty"` //Satış fiyatı
	StockQuantity int                `json:"stock_quantity,omitempty" bson:"stock_quantity,omitempty"`
	Barcode       uint               `json:"barcode,,omitempty" bson:"barcode,omitempty"`  //Barcode numarası
	StoreID       uint               `json:"store_id,omitempty" bson:"store_id,omitempty"` //StoreID
	Categories    []Category         `json:"categories,omitempty" bson:"categories,omitempty"`
	Entry_Price   float64            `json:"entry_price,omitempty" bson:"entry_price,omitempty"` //Alış fiyatı
	Tax           float64            `json:"kdv,omitempty" bson:"kdv,omitempty"`                 //Vergi oranı
	SalesQty      int                `json:"salesqty,omitempty" bson:"salesqty,omitempty"`       //Order başına satış adedi
}
