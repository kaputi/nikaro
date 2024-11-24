package server

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Age    int                `json:"age" bson:"age"`
	Gender bool               `json:"gender" bson:"gender"`
}
