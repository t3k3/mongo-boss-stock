package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repair struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	Tel             string             `json:"tel,omitempty" bson:"tel,omitempty"`
	Problem         string             `json:"problem,omitempty" bson:"problem,omitempty"` //Problem fiyatı
	Status          string             `json:"status,omitempty" bson:"status,omitempty"`
	Notes           string             `json:"notes,,omitempty" bson:"notes,omitempty"` //Notlar
	Estimated_price float64            `json:"estimated_price,omitempty" bson:"estimated_price,omitempty"`
	Brand           string             `json:"brand" bson:"brand"`
	Device_model    string             `json:"device_model,omitempty" bson:"device_model,omitempty"` //Alış fiyatı
	Color           string             `json:"color,omitempty" bson:"color,omitempty"`               //Vergi oranı
	Diagnosis       string             `json:"diagnosis,omitempty" bson:"diagnosis,omitempty"`
	Sms             bool               `json:"sms,omitempty" bson:"sms,omitempty"` //Problem
}
