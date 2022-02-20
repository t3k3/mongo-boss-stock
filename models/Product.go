package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt     time.Time          `json:"product_created_date,omitempty" bson:"product_created_date,omitempty"`
	DeletedAt     time.Time          `json:"product_deleted_date,omitempty" bson:"product_deleted_date,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	Detail        string             `json:"detail,omitempty" bson:"detail,omitempty"`
	Slug          string             `json:"slug,omitempty" bson:"slug,omitempty"`
	Price         int                `json:"price,omitempty" bson:"price,omitempty"` //Satış fiyatı
	StockQuantity int                `json:"stock_quantity,omitempty" bson:"stock_quantity,omitempty"`
	Barcode       string             `json:"barcode,omitempty" bson:"barcode,omitempty"` //Barcode numarası
	Category      string             `json:"category,omitempty" bson:"category,omitempty"`
	CartQuantity  int                `json:"cartQuantity,omitempty" bson:"cartQuantity,omitempty"`
	ImageURL      string             `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	SalesQty      int                `json:"salesqty,omitempty" bson:"salesqty,omitempty"` //Satış per order
	Status        bool               `json:"status,omitempty" bson:"status,omitempty"`
	Entry_Price   int                `json:"entry_price,omitempty" bson:"entry_price,omitempty"` //Alış fiyatı
	Tedarikci     string             `json:"tedarikci,omitempty" bson:"tedarikci,omitempty"`

	//Categories    []Category         `json:"categories,omitempty" bson:"categories,omitempty"`
	//StoreID       uint               `json:"store_id,omitempty" bson:"store_id,omitempty"` //StoreID

	//Tax           float64            `json:"kdv,omitempty" bson:"kdv,omitempty"`                 //Vergi oranı
	//Brand         string             `json:"brand,omitempty" bson:"brand,omitempty"`
	//Device_model  string             `json:"device_model,omitempty" bson:"device_model,omitempty"` //Alış fiyatı
	//Guaranty      string             `json:"device_guaranty,omitempty" bson:"device_guaranty,omitempty"`
	//Color         string             `json:"device_color,omitempty" bson:"device_color,omitempty"`
	//Capacity      string             `json:"device_capacity,omitempty" bson:"device_capacity,omitempty"`
}
