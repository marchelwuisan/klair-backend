package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Wallet struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserId     string             `json:"user_id,omitempty"`
	Name       string             `json:"name,omitempty" validate:"required"`
	Type       string             `json:"type,omitempty" validate:"required"`
	Currency   string             `json:"currency,omitempty" validate:"required"`
	Balance    int                `json:"balance,omitempty" validate:"required"`
	Records    []Record           `json:"record,omitempty"`
	CreatedAt  int64              `json:"createdat,omitempty"`
	UpdatedAt  int64              `json:"updatedat,omitempty"`
	DeletedAt  int64              `json:"deletedat,omitempty"`
	DataStatus int                `json:"dataStatus,omitempty"`
}
