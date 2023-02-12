package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Record struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserId       string             `json:"user_id,omitempty"`
	WalletId     string             `json:"wallet_id,omitempty"`
	DebtorId     string             `json:"debtor,omitempty"`
	WalletFromId string             `json:"wallet_from_id,omitempty"`
	WalletToId   string             `json:"wallet_to_id,omitempty"`
	Amount       int                `json:"amount,omitempty" validate:"required"`
	Type         string             `json:"type,omitempty" validate:"required"`
	Category     string             `json:"category,omitempty" validate:"required"`
	Note         string             `json:"note,omitempty"`
	Attachment   string             `json:"attachment,omitempty"`
	CreatedAt    int64              `json:"createdat,omitempty"`
	UpdatedAt    int64              `json:"updatedat,omitempty"`
	DeletedAt    int64              `json:"deletedat,omitempty"`
	DataStatus   int                `json:"dataStatus,omitempty"`
}
