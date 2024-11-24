package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/kaputi/nikaro/internal/utils"
)

var loadedEnv = false

func SetupEnv() {
	err := godotenv.Load()
	utils.MustErr(err)
	loadedEnv = true
}

func EnvMongoURI() string {
	utils.MustOk(loadedEnv, "Env not loaded")
	user, ok := os.LookupEnv("MONGO_USER")
	utils.MustOk(ok, "MONGO_USER not found")
	password, ok := os.LookupEnv("MONGO_USER_PASSWORD")
	utils.MustOk(ok, "MONGO_USER_PASSWORD not found")
	port, ok := os.LookupEnv("MONGO_PORT")
	utils.MustOk(ok, "MONGO_PORT not found")

	mongoUrl := fmt.Sprintf("mongodb://%s:%s@localhost:%s", user, password, port)

	return mongoUrl
}

func EnvServerPort() string {
	utils.MustOk(loadedEnv, "Env not loaded")
	port, ok := os.LookupEnv("PORT")
	utils.MustOk(ok, "PORT not found")
	return port
}
