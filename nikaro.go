package main

import (
	"github.com/kaputi/nikaro/internal/configs"
	"github.com/kaputi/nikaro/internal/database"
	"github.com/kaputi/nikaro/internal/server"
)

func main() {
	configs.SetupEnv()

	client, ctx, cancel := database.ConnectDB()

	// TODO: SERVER CODE GOES HERE
	server.Serve()

	database.CloseConnection(client, ctx, cancel)
}
