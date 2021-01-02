package main

import (
	"fmt"
	"net/http"
	"os"
	"soci-avatar-cdn/config"
	"soci-avatar-cdn/route"
)

func setupRoutes(settings *config.Config) {
	http.Handle("/", http.FileServer(http.Dir("./files")))
	http.HandleFunc("/upload", route.UploadFile)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = settings.Port
		if port == "" {
			port = "4202"
		}
	}

	fmt.Printf("Listening on %v\n", port)
	http.ListenAndServe(":"+port, nil)
}

func main() {
	var settings config.Config
	// parse the config file
	if err := config.ParseJSONFile("./config.json", &settings); err != nil {
		panic(err)
	}
	// validate the config file
	if err := settings.Validate(); err != nil {
		panic(err)
	}

	fmt.Println("Starting avatar encoding server...")
	setupRoutes(&settings)
}
