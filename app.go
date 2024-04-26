package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Quality int

const (
	HIGH Quality = iota
	MEDIUM
	LOW
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

	// monitor the App resources on port 8081
	go app.Monitor()

	app.MuxRouter = mux.NewRouter().StrictSlash(true)
	app.MuxRouter.HandleFunc("/api/v1/terminate", app.Terminate)
	app.MuxRouter.HandleFunc("/api/v1/transcode/{quality}", app.Transcode)
	log.Fatal(http.ListenAndServe(":8080", app.MuxRouter))
}

// Terminate - cleanly close all go routines and recover resources
func (app *App) Terminate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Terminating server")
}

func (app *App) ErrorResponse(w http.ResponseWriter, code int, item string) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("{\"error\":%s}", item)))
}

func (app *App) ValidResponse(w http.ResponseWriter, item string) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("{\"id\":%s}", item)))
}
