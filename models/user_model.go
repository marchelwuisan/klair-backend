package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Login struct {
	Email string `json:"email,omitempty" validate:"required"`
}

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Email         string             `json:"email,omitempty" validate:"required"`
	Username      string             `json:"username,omitempty" validate:"required"`
	FirstName     string             `json:"firstname,omitempty" validate:"required"`
	LastName      string             `json:"lastname,omitempty" validate:"required"`
	Currency      string             `json:"currency,omitempty" validate:"required"`
	Phone         string             `json:"phone" validate:"required"`
	Password      string             `json:"password,omitempty" validate:"required"`
	FirebaseUid   string             `json:"firebaseuid,omitempty"`
	PushToken     string             `json:"pushtoken,omitempty"`
	Token         string             `json:"token,omitempty"`
	RefreshToken  string             `json:"refreshtoken,omitempty"`
	Gender        string             `json:"gender,omitempty"`
	IsFirstSignIn int                `json:"isfirstsignin,omitempty"`
	DateOfBirth   int64              `json:"dateofbirth,omitempty"`
	CreatedAt     int64              `json:"createdat,omitempty"`
	UpdatedAt     int64              `json:"updatedat,omitempty"`
	DeletedAt     int64              `json:"deletedat,omitempty"`
	DataStatus    int                `json:"dataStatus,omitempty"`
	// User_id      string             `json:"user_id,omitempty"`
}
