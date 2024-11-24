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
	"github.com/kaputi/nikaro/internal/configs"
	"github.com/kaputi/nikaro/internal/database"
	"github.com/kaputi/nikaro/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTest struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Pass string             `json:"pass" bson:"pass"`
}

func Serve() {

	port := configs.EnvServerPort()

	server := &http.Server{Addr: ":" + port, Handler: service()}

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
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)

		utils.FatalErr(err)

		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func service() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Go then, there are other worlds than these.\n"))
		utils.LogErr(err)
	})

	// TEST REGISTER
	r.Post("/api/v1/auth", func(w http.ResponseWriter, r *http.Request) {

		user := UserTest{
			ID: primitive.NewObjectID(),
		}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		usersCollection := database.GetCollection("Users")
		_, err = usersCollection.InsertOne(context.Background(), user)

		if err != nil {
			fmt.Println("Error inserting user: ", err)
		}

		_, err = w.Write([]byte("Auth endpoint\n"))
		utils.LogErr(err)
	})

	return r
}
