package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserId string             `json:"user_id,omitempty"`
	Name   string             `json:"name,omitempty" validate:"required"`
}
