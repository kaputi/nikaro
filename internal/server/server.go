package server

import (
	"context"
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
	"github.com/kaputi/nikaro/internal/modules/user"
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
	router.Handle("/*", http.FileServer(http.Dir("./public")))

	return router
}
