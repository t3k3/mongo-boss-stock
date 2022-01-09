package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Category_name string             `json:"category_name,omitempty" bson:"category_name,omitempty"`
	Sub_category  uint               `json:"sub_category,omitempty" bson:"sub_category,omitempty"`
}
