package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaputi/nikaro/internal/database"
	"github.com/kaputi/nikaro/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	UserName string             `json:"userName" bson:"username"`
	Password string             `json:"password" bson:"password"`
}

type UserPayload struct {
	Id       primitive.ObjectID `json:"id"`
	UserName string             `json:"userName"`
}

type loginModel struct {
	UserName string `json:"userName" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		collection: database.GetCollection("users"),
	}
}

func (ur *UserRepo) Routes() chi.Router {
	r := chi.NewRouter()

	// TODO: list needs authorization middleware
	r.Get("/", ur.List)

	r.Post("/", ur.Create)

	r.Post("/login", ur.Login)

	// TODO: user specific routes needs authentication or authorization (admin has access to all user) middleware
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", ur.Get)
		r.Put("/", ur.Update)
		r.Delete("/", ur.Delete)
	})

	return r
}

func (ur *UserRepo) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List users"))
}

func (ur *UserRepo) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get user"))
}

func (ur *UserRepo) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation
	user := User{
		Id: primitive.NewObjectID(),
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pass, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = pass

	_, err = ur.collection.InsertOne(r.Context(), user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: set cookie with jwt

	w.Header().Set("Content-Type", "application/json")
	responseUser := UserPayload{
		Id:       user.Id,
		UserName: user.UserName,
	}

	err = json.NewEncoder(w).Encode(responseUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ur *UserRepo) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation
	reqData := loginModel{}

	err := json.NewDecoder(r.Body).Decode(&reqData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := User{}

	err = ur.collection.FindOne(r.Context(), bson.M{"username": reqData.UserName}).Decode(&user)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	if utils.VerifyPassword(reqData.Password, user.Password) {
		// TODO: set cookie with jwt
		w.Header().Set("Content-Type", "application/json")
		responseUser := UserPayload{
			Id:       user.Id,
			UserName: user.UserName,
		}

		err = json.NewEncoder(w).Encode(responseUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

}

func (ur *UserRepo) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update user"))
}

func (ur *UserRepo) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
