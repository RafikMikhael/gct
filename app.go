package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	NumberOfRenditions  = 5
	NumberOfQualityLvls = 3
)

type App struct {
	MuxRouter *mux.Router
	Jobs      map[string]*Job
	// look-up tables
	bitRate   [NumberOfQualityLvls][NumberOfRenditions]int //[quality][renditionIdx]
	horizW    [NumberOfRenditions]int                      //rendition target width in pixels
	vertH     [NumberOfRenditions]int                      //rendition target height in pixels
	sleepTime [NumberOfRenditions]int                      //in seconds

}

// Initialize - initialize App fields and allocate all needed memory
func (app *App) Initialize() {
	app.bitRate = [NumberOfQualityLvls][NumberOfRenditions]int{
		{160, 360, 1930, 4080, 7000},
		{145, 300, 1600, 3400, 5800},
		{120, 280, 1400, 3080, 4500},
	}
	app.horizW = [NumberOfRenditions]int{640, 768, 960, 1280, 1920}
	app.vertH = [NumberOfRenditions]int{360, 432, 540, 720, 1080}
	app.sleepTime = [NumberOfRenditions]int{10, 20, 30, 40, 50}
	app.Jobs = make(map[string]*Job)
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
	app.MuxRouter.HandleFunc("/api/v1/probe/{hash}", app.ProbeHash)
	log.Fatal(http.ListenAndServe(":8080", app.MuxRouter))
}

// Terminate - cleanly close all go routines and recover resources
func (app *App) Terminate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Terminating server")
}

func (app *App) ErrorResponse(w http.ResponseWriter, code int, key, value string) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("{\"%s\":%s}", key, value)))
}
