package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

const (
	NumberOfRenditions  = 5
	NumberOfQualityLvls = 3
)

type App struct {
	MuxRouter *mux.Router
	mu        sync.Mutex
	Jobs      map[string]*Job
	bStopped  bool
	port      string

	// stuff needed for graceful shutdown of server
	Cancel      context.CancelFunc
	Ctx         context.Context
	StopSignals chan os.Signal

	// look-up tables
	bitRate   [NumberOfQualityLvls][NumberOfRenditions]int //[quality][renditionIdx]
	horizW    [NumberOfRenditions]int                      //rendition target width in pixels
	vertH     [NumberOfRenditions]int                      //rendition target height in pixels
	sleepTime [NumberOfRenditions]int                      //in seconds
}

// Initialize - initialize App fields and allocate all needed memory
func (app *App) Initialize(portNum *int) {
	app.bitRate = [NumberOfQualityLvls][NumberOfRenditions]int{
		{160, 360, 1930, 4080, 7000},
		{145, 300, 1600, 3400, 5800},
		{120, 280, 1400, 3080, 4500},
	}
	app.horizW = [NumberOfRenditions]int{640, 768, 960, 1280, 1920}
	app.vertH = [NumberOfRenditions]int{360, 432, 540, 720, 1080}
	app.sleepTime = [NumberOfRenditions]int{10, 20, 30, 40, 50}
	app.Jobs = make(map[string]*Job)
	app.port = ":" + strconv.Itoa(*portNum)

	app.StopSignals = make(chan os.Signal, 1)
	signal.Notify(app.StopSignals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

}

// Run - run the application (main go routine running forever)
func (app *App) Run() {
	// graceful shutdown can be triggered by /terminate or by ctrl+c
	// https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97
	app.Ctx, app.Cancel = context.WithCancel(context.Background())
	go func() {
		<-app.StopSignals
		app.Cancel()
	}()

	app.MuxRouter = mux.NewRouter().StrictSlash(true)
	app.MuxRouter.HandleFunc("/api/v1/terminate", app.Terminate)
	app.MuxRouter.HandleFunc("/api/v1/job/{quality}", app.TriggerJobs)
	app.MuxRouter.HandleFunc("/api/v1/probe/{hash}", app.ProbeHash)

	server := http.Server{
		Addr:    app.port,
		Handler: app.MuxRouter,
	}
	// we can add any middleware here like corsHandler or LoggingHandler
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	// monitor the App resources on port 8081
	go app.Monitor()

	log.Printf("transcode server starting up")
	<-app.Ctx.Done()
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("transcode server shutting down")
}

// Terminate - cleanly close all go routines and recover resources
func (app *App) Terminate(w http.ResponseWriter, r *http.Request) {
	// verify the verb used
	if r.Method != "GET" {
		// 405
		app.JsonHttpResponse(w, http.StatusMethodNotAllowed, "error", r.Method)
		return
	}
	app.bStopped = true
	app.JsonHttpResponse(w, http.StatusOK, "termination", "started")
}

func (app *App) JsonHttpResponse(w http.ResponseWriter, code int, key, value string) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("{\"%s\":%s}", key, value)))
}
