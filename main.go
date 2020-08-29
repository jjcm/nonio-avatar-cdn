package main

import (
	"fmt"
	"net/http"
	"os"
	"soci-cdn/route"
)

func setupRoutes() {
	http.HandleFunc("/upload", route.UploadFile)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4202"
	}

	http.ListenAndServe(":"+port, nil)
}

func main() {
	fmt.Println("Starting media encoding server")
	setupRoutes()
}
