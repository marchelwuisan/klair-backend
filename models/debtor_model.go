package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Debtor struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `json:"name,omitempty" validate:"required"`
	UserId            string             `json:"user_id,omitempty"`
	PayableBalance    int                `json:"payable_balance,omitempty"`
	RecievableBalance int                `json:"recievable_balance,omitempty"`
	Records           []Record           `json:"record,omitempty"`
	CreatedAt         int64              `json:"createdat,omitempty"`
	UpdatedAt         int64              `json:"updatedat,omitempty"`
	DeletedAt         int64              `json:"deletedat,omitempty"`
	DataStatus        int                `json:"dataStatus,omitempty"`
}
