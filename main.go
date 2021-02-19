package main

import (
	"log"
	"speedСontrol/services/server"
)

func main() {
	httpServer := server.GetHttpServer()
	log.Printf("Server is listening on address %s ... \n", httpServer.Addr)

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
