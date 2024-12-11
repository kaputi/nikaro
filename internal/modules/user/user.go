package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kaputi/nikaro/internal/auth"
	"github.com/kaputi/nikaro/internal/database"
	"github.com/kaputi/nikaro/internal/res"
	"github.com/kaputi/nikaro/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Role     string             `json:"role" bson:"role"`
	Password string             `json:"password" bson:"password"`
}

type UserPayload struct {
	Id       primitive.ObjectID `json:"id"`
	UserName string             `json:"username"`
}

type UserRepo struct {
	collection *mongo.Collection
	auth       *auth.Authorization
}

func NewUserRepo(authorization *auth.Authorization) *UserRepo {

	collection := database.GetCollection("users")

	// TODO:
	// collection.Indexes().CreateOne(
	// 	context.Background(),
	// 	mongo.IndexModel{
	// 		Keys:    bson.D{{Key: "username", Value: 1}},
	// 		Options: options.Index().SetUnique(true),
	// 	},
	// )

	return &UserRepo{
		collection: collection,
		auth:       authorization,
	}
}

func (ur *UserRepo) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", ur.Create)

	r.Post("/login", ur.Login)

	r.Group(func(r chi.Router) {
		r.Use(ur.auth.VerifyToken("jwt"))
		r.Route("/", func(r chi.Router) {
			r.Get("/", ur.Get)
			r.Put("/", ur.Update)
			r.Delete("/", ur.Delete)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(ur.auth.VerifyToken("refresh_token"))
		r.Get("/refresh", ur.RefreshToken)
	})

	r.Group(func(r chi.Router) {
		r.Use(ur.auth.VerifyToken("jwt"))
		r.Use(ur.auth.AuthorizeAdmin())
		r.Get("/list", ur.List)
	})

	return r
}

func (ur *UserRepo) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List users"))
}

func (ur *UserRepo) Get(w http.ResponseWriter, r *http.Request) {
	_, claims := ur.auth.GetTokenFromContext(r.Context())

	// if claims.Subject != "" {
	// 	res.Success(w, claims.Subject)
	//    return
	// }
	id, _ := primitive.ObjectIDFromHex(claims.Subject)

	user := UserPayload{}

	err := ur.collection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Fail(w, "user not found", http.StatusNotFound)
			return
		}
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Success(w, user)
	// responseUser := UserPayload{
	// 	Id:       user.Id,
	// 	UserName: user.Username,
	// }

}

func (ur *UserRepo) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation

	user := User{
		Role: "admin", // TODO: this is hardcoded for now
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		res.BadRequest(w, err.Error())
		return
	}
	user.Id = primitive.NewObjectID()

	inDb := ur.collection.FindOne(r.Context(), bson.M{"username": user.Username})
	if inDb.Err() == nil {
		res.Fail(w, "user already exists", http.StatusConflict)
		return
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = hashed

	_, err = ur.collection.InsertOne(r.Context(), user)

	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ur.setBothTokens(w, r, user)
	if err != nil {
		return
	}

	responseUser := UserPayload{
		Id:       user.Id,
		UserName: user.Username,
	}

	res.Success(w, responseUser)
}

func (ur *UserRepo) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation

	reqData := User{}

	err := json.NewDecoder(r.Body).Decode(&reqData)

	if err != nil {
		res.BadRequest(w, err.Error())
		return
	}

	user := User{}

	err = ur.collection.FindOne(r.Context(), bson.M{"username": reqData.Username}).Decode(&user)
	if err != nil {
		res.Fail(w, "user not found", http.StatusUnauthorized)
		return
	}

	if utils.VerifyPassword(reqData.Password, user.Password) {
		err = ur.setBothTokens(w, r, user)
		if err != nil {
			res.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseUser := UserPayload{
			Id:       user.Id,
			UserName: user.Username,
		}

		res.Success(w, responseUser)

		return
	}

	res.Fail(w, "invalid password", http.StatusUnauthorized)
}

func (ur *UserRepo) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update user"))
}

func (ur *UserRepo) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}

func (ur *UserRepo) RefreshToken(w http.ResponseWriter, r *http.Request) {
	_, claims := ur.auth.GetTokenFromContext(r.Context())
	id, err := primitive.ObjectIDFromHex(claims.Subject)

	if err != nil {
		res.BadRequest(w, "Invalid user id in token")
		return
	}

	user := User{
		Id:   id,
		Role: claims.Role,
	}

	err = ur.setBothTokens(w, r, user)
	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Success(w, "Token refreshed")
}

func (ur *UserRepo) setTokenCookie(w http.ResponseWriter, r *http.Request, user User) (string, error) {
	token, err := ur.auth.GenerateToken(user.Id.Hex(), user.Role, time.Minute*5)
	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	ur.auth.SetTokenToCookie(w, "jwt", token, "/", time.Minute*5)
	return token, nil
}

func (ur *UserRepo) setRefreshTokenCookie(w http.ResponseWriter, r *http.Request, user User) (string, error) {
	token, err := ur.auth.GenerateToken(user.Id.Hex(), user.Role, time.Hour*24*7)
	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	ur.auth.SetTokenToCookie(w, "refresh_token", token, utils.ApiRoute("auth/refresh"), time.Minute*5)
	return token, nil
}

func (ur *UserRepo) setBothTokens(w http.ResponseWriter, r *http.Request, user User) error {
	_, err := ur.setTokenCookie(w, r, user)
	if err != nil {
		return err
	}
	_, err = ur.setRefreshTokenCookie(w, r, user)
	if err != nil {
		return err
	}
	return nil
}
