package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kaputi/nikaro/internal/configs"
	"github.com/kaputi/nikaro/internal/database"
	"github.com/kaputi/nikaro/internal/modules/drawings"
	"github.com/kaputi/nikaro/internal/modules/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestServer struct {
	httpServer *http.Server

	userStore *user.UserRepo
	// drawingStore *DrawingStore
	// collabStore *CollabStore
}

func CreateRestServer() *RestServer {
	return &RestServer{
		userStore: user.NewUserRepo(),
	}
}

func (rs *RestServer) Start() {

	port := configs.EnvServerPort()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: rs.Routes(),
	}

	rs.httpServer = server

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			defer shutdownStopCtx()
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("[Debug] graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)

		if err != nil {
			log.Fatalf("[Debug] %s", err)
		}

		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func apiRoute(route string) string {
	apiV, ok := os.LookupEnv("API_VERSION")
	if !ok {
		apiV = "/api/v1/"
	}

	return fmt.Sprintf("%s%s", apiV, route)
}

func (rs *RestServer) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(middleware.Throttle(200))

	// TODO: check for this on prod
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders: []string{"Link"},
    AllowCredentials: false,
    MaxAge: 300,
	}))

	// Health check
	router.Get(apiRoute("yougood"), func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("I'm good!"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	router.Mount(apiRoute("auth"), rs.userStore.Routes())

	// staitc
	router.Handle("/*", http.FileServer(http.Dir(configs.EnvFrontEndDir())))

	router.Get("/api/v1/drawings", func(w http.ResponseWriter, r *http.Request) {
		collection := database.GetCollection("drawings")
		allDrawings := drawings.Drawing{}

		obectId, _ := primitive.ObjectIDFromHex("67462822fa632306906a5d96")
		err := collection.FindOne(r.Context(), bson.M{"_id": obectId}).Decode(&allDrawings)

		if err != nil {
			//TODO: handle error
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(allDrawings.Elements)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})

	router.Post("/drawings", func(w http.ResponseWriter, r *http.Request) {
		dr := drawings.Drawing{
			Id:       primitive.NewObjectID(),
			Elements: []drawings.ExcalidrawElement{},
		}

		err := json.NewDecoder(r.Body).Decode(&dr.Elements)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		collection := database.GetCollection("drawings")
		_, err = collection.InsertOne(r.Context(), dr)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(dr.Elements)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})

	return router
}
