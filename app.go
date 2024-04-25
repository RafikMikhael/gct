package main

import (
	"log"
)

type App struct {
}

// Initialize - initialize App fields and allocate all needed memory
func (app *App) Initialize() {
}

// Run - run the application (main go routine running forever)
func (app *App) Run() {
	log.Print("long-stack server starting up")
	log.Print("")
	// Intentionally not worrying about graceful shutdown of the server
	// e.g. https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97
	defer log.Print("long-stack server shutting down")

	// monitor the App resources on port 8081
	go app.Monitor()
}
