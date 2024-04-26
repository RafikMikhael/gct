package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	MuxRouter *mux.Router
	bitRate   [3][5]int //[quality][renditionIdx]
	horizW    [5]int    //rendition target width in pixels
	vertH     [5]int    //rendition target height in pixels
	sleepTime [5]int    //in seconds
}

// Initialize - initialize App fields and allocate all needed memory
func (app *App) Initialize() {
	app.bitRate = [3][5]int{
		{160, 360, 1930, 4080, 7000},
		{145, 300, 1600, 3400, 5800},
		{120, 280, 1400, 3080, 4500},
	}
	app.horizW = [5]int{640, 768, 960, 1280, 1920}
	app.vertH = [5]int{360, 432, 540, 720, 1080}
	app.sleepTime = [5]int{1, 2, 3, 4, 5}
}

// Run - run the application (main go routine running forever)
func (app *App) Run() {
	log.Print("transcode server starting up")
	defer log.Print("transcode server shutting down")

	// monitor the App resources on port 8081
	go app.Monitor()

	app.MuxRouter = mux.NewRouter().StrictSlash(true)
	app.MuxRouter.HandleFunc("/api/v1/terminate", app.Terminate)
	app.MuxRouter.HandleFunc("/api/v1/job/{quality}", app.TriggerJobs)
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
