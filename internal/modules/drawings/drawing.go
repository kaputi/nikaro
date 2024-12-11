package drawings

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaputi/nikaro/internal/auth"
	"github.com/kaputi/nikaro/internal/database"
	"github.com/kaputi/nikaro/internal/res"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Drawing struct {
	Id       primitive.ObjectID  `json:"_id" bson:"_id"`
	UserId   primitive.ObjectID  `json:"userId" bson:"userId"`
	Name     string              `json:"name" bson:"name"`
	Elements []ExcalidrawElement `json:"drawings" bson:"drawings"`
}

type DrawingsRepo struct {
	collection *mongo.Collection
	auth       *auth.Authorization
}

func NewDrawingsRepo(authorization *auth.Authorization) *DrawingsRepo {
	return &DrawingsRepo{
		collection: database.GetCollection("drawings"),
		auth:       authorization,
	}
}

func (dr *DrawingsRepo) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(dr.auth.VerifyToken("jwt"))
		r.Get("/", dr.List)
		r.Post("/", dr.Create)
		r.Route("/{drawingId}", func(r chi.Router) {
			r.Get("/", dr.Get)
			r.Put("/", dr.Update)
			r.Delete("/", dr.Delete)
		})
	})

	return r
}

func (dr *DrawingsRepo) List(w http.ResponseWriter, r *http.Request) {
	res.Success(w, "List")
}

func (dr *DrawingsRepo) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation
	drawing := Drawing{}

	err := json.NewDecoder(r.Body).Decode(&drawing)
	if err != nil {
		res.BadRequest(w, err.Error())
		return
	}
	drawing.Id = primitive.NewObjectID()

	_, claims := dr.auth.GetTokenFromContext(r.Context())
	id, err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil {
		res.BadRequest(w, "Invalid user id in token")
		return
	}

	drawing.UserId = id

	inDb := dr.collection.FindOne(r.Context(), bson.M{"name": drawing.Name, "userId": drawing.UserId})
	if inDb.Err() == nil {
		res.Fail(w, "Drawing with this name already exists", http.StatusConflict)
		return
	}

	_, err = dr.collection.InsertOne(r.Context(), drawing)

	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Success(w, drawing)
}

func (dr *DrawingsRepo) Get(w http.ResponseWriter, r *http.Request) {
	drawingId := chi.URLParam(r, "drawingId")
	Id, err := primitive.ObjectIDFromHex(drawingId)
	if err != nil {
		res.BadRequest(w, "invalid drawing id")
		return
	}

	_, claims := dr.auth.GetTokenFromContext(r.Context())
	userId, _ := primitive.ObjectIDFromHex(claims.Subject)

	drawing := Drawing{Id: Id, UserId: userId}

	findBson := bson.M{"_id": drawing.Id, "userId": drawing.UserId}

	err = dr.collection.FindOne(r.Context(), findBson).Decode(&drawing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Fail(w, "Drawing not found", http.StatusNotFound)
			return
		}
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(drawing)
	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dr *DrawingsRepo) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update"))
}

func (dr *DrawingsRepo) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete"))
}
