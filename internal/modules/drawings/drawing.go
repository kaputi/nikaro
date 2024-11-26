package drawings

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaputi/nikaro/internal/database"
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
}

func NewDrawingsRepo() *DrawingsRepo {
	return &DrawingsRepo{
		collection: database.GetCollection("drawings"),
	}
}

func (dr *DrawingsRepo) Routes() chi.Router {
	r := chi.NewRouter()
	// TODO: all this needs auth middleware
	r.Get("/", dr.List)
	r.Post("/", dr.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", dr.Get)
		r.Put("/", dr.Update)
		r.Delete("/", dr.Delete)
	})
	return r
}

func (dr *DrawingsRepo) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List"))
}

func (dr *DrawingsRepo) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation
	drawing := Drawing{
		Id: primitive.NewObjectID(),
	}

	err := json.NewDecoder(r.Body).Decode(&drawing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = dr.collection.InsertOne(r.Context(), drawing)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(drawing)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dr *DrawingsRepo) Get(w http.ResponseWriter, r *http.Request) {
	drawing := Drawing{
		Id: primitive.NewObjectID(),
	}

	err := json.NewDecoder(r.Body).Decode(&drawing)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: check if this works  without primitive.ObjectIdFromHex
	findBson := bson.M{"_id": drawing.Id, "userId": drawing.UserId}

	err = dr.collection.FindOne(r.Context(), findBson).Decode(&drawing)
	if err != nil {
		// TODO:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(drawing)

}

func (dr *DrawingsRepo) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update"))
}

func (dr *DrawingsRepo) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete"))
}
