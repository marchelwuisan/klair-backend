package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Goal struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	WalletId   string             `json:"wallet_id,omitempty"`
	UserId     string             `json:"user_id,omitempty"`
	Name       string             `json:"name,omitempty" validate:"required"`
	Amount     string             `json:"amount,omitempty" validate:"required"`
	Category   string             `json:"category,omitempty"`
	CreatedAt  int64              `json:"createdat,omitempty"`
	UpdatedAt  int64              `json:"updatedat,omitempty"`
	DeletedAt  int64              `json:"deletedat,omitempty"`
	DataStatus int                `json:"dataStatus,omitempty"`
}
