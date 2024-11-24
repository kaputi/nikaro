package main

import (
	"fmt"

	"github.com/kaputi/nikaro/internal/server"
)

func main() {
	client, ctx, cancel := server.ConectMongoDb()
	fmt.Println("client", client)
	fmt.Println("ctx", ctx)
	fmt.Println("cancel", cancel)

	server.CloseConnection(client, ctx, cancel)
}
