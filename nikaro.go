package main

import (
	"github.com/kaputi/nikaro/internal/configs"
	"github.com/kaputi/nikaro/internal/database"
	"github.com/kaputi/nikaro/internal/server"
)

func main() {
	configs.SetupEnv()

	client, ctx, cancel := database.ConnectDB()

	sr := server.CreateRestServer()
	sr.Routes()
	sr.Start()

	database.CloseConnection(client, ctx, cancel)
}
