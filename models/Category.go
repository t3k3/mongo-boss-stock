package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Category_name string             `json:"category_name" bson:"category_name,omitempty"`
	Sub_category  uint               `json:"sub_category" bson:"sub_category,omitempty"`
}
