package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Address  string             `json:"adres,omitempty" bson:"adres,omitempty"`
	City     string             `json:"city,omitempty" bson:"city,omitempty"`
	Region   string             `json:"region,omitempty" bson:"region,omitempty"`
	Tel      string             `json:"tel,omitempty" bson:"tel,omitempty"`
	Mail     string             `json:"mail,omitempty" bson:"mail,omitempty"`
	IsActive bool               `json:"isactive,omitempty" bson:"isactive,omitempty"`
	Balance  float64            `json:"balance,omitempty" bson:"balance,omitempty"`
}
