package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintf(w, "Error in reading post body %s", err)
		return
	}
	collection, client, context, cancel := SetupMongoDB()
	defer CloseConnection(client, context, cancel)
	user.Id = primitive.NewObjectID()
	result, err := collection.InsertOne(context, user)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error in Creating user %v", err))
		return
	}
	fmt.Printf("Inserted user %v", result.InsertedID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

// func (uh UserHandler) FindUsers(w http.ResponseWriter, r *http.Request) {
// 	collection, client, context, cancel := SetupMongoDB()
// 	u := make([]dto.User, 0, 10)
// 	filter := bson.D{}
// 	cursor, err := collection.Find(context, filter)
// 	if err == mongo.ErrNoDocuments {
// 		fmt.Println("No document found")
// 	} else if err != nil {
// 		fmt.Printf("Error in mongo %v", err)
// 	}
// 	if err != nil {
// 		fmt.Fprintf(w, fmt.Sprintf("Error in Finding user %v", err))
// 	}
// 	defer CloseConnection(client, context, cancel)
// 	defer cursor.Close(context)
// 	for cursor.Next(context) {
// 		var result User
// 		err := cursor.Decode(&result)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		u = append(u, result)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(u); err != nil {
// 		fmt.Fprintf(w, fmt.Sprintf("Cannot parse %v", err))
// 	}
// }
