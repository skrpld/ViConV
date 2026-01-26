package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	UserId                 bson.ObjectID `bson:"_id"` //TODO: bson.ObjectID -> мб новый тип в consts
	Email                  string        `bson:"email"`
	PasswordHash           string        `bson:"password_hash"`
	RefreshToken           string        `bson:"refresh_token"`
	RefreshTokenExpiryTime time.Time     `bson:"refresh_token_expiry_time"`
}
