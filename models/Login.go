package models

type Login struct {
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Pass  string `json:"password,omitempty" bson:"password,omitempty"`
}
