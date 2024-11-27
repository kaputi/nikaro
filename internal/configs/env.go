package configs

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/kaputi/nikaro/internal/utils"
)

func SetupEnv() {
	err := godotenv.Load()
	utils.MustErr(err)

	mongoUser, ok := os.LookupEnv("MONGO_USER")
	utils.MustOk(ok, "MONGO_USER not found")
	utils.MustOk(mongoUser != "", "MONGO_USER")

	mongoPassword, ok := os.LookupEnv("MONGO_USER_PASSWORD")
	utils.MustOk(ok, "MONGO_USER_PASSWORD not found")
	utils.MustOk(mongoPassword != "", "MONGO_USER_PASSWORD")

	mongoPort, ok := os.LookupEnv("MONGO_PORT")
	utils.MustOk(ok, "MONGO_PORT not found")
	utils.MustOk(mongoPort != "", "MONGO_PORT")

	port, ok := os.LookupEnv("PORT")
	utils.MustOk(ok, "PORT not found")
	utils.MustOk(port != "", "PORT")

	frontEndDir, ok := os.LookupEnv("FRONT_END_BUILD_DIR")
	utils.MustOk(ok, "FRONT_END_BUILD_DIR not found")
	utils.MustOk(frontEndDir != "", "FRONT_END_BUILD_DIR")

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	utils.MustOk(ok, "JWT_SECRET not found")
	utils.MustOk(jwtSecret != "", "JWT_SECRET")
}
