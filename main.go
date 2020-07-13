package main

import (
	"fmt"
	"net/http"
	"soci-cdn/route"
)

func setupRoutes() {
	http.HandleFunc("/upload", route.UploadFile)
	http.ListenAndServe(":8082", nil)
}

func main() {
	fmt.Println("Starting media encoding server")
	setupRoutes()
}
