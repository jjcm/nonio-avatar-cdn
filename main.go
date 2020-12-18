package main

import (
	"fmt"
	"net/http"
	"os"
	"soci-avatar-cdn/route"
)

func setupRoutes() {
	http.HandleFunc("/upload", route.UploadFile)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4202"
	}

	fmt.Printf("Listening on %v\n", port)
	http.ListenAndServe(":"+port, nil)
}

func main() {
	fmt.Println("Starting avatar encoding server...")
	setupRoutes()
}
