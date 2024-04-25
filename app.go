package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	MuxRouter *mux.Router
}

// Initialize - initialize App fields and allocate all needed memory
func (app *App) Initialize() {
}

// Run - run the application (main go routine running forever)
func (app *App) Run() {
	log.Print("transcode server starting up")
	defer log.Print("transcode server shutting down")
	// Intentionally not worrying about graceful shutdown of the server
	// e.g. https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97

	// monitor the App resources on port 8081
	go app.Monitor()

	app.MuxRouter = mux.NewRouter().StrictSlash(true)
	app.MuxRouter.HandleFunc("/terminate", app.Terminate)
	log.Fatal(http.ListenAndServe(":8080", app.MuxRouter))
}

// Terminate - cleanly close all go routines and recover resources
func (app *App) Terminate(w http.ResponseWriter, r *http.Request) {

}
