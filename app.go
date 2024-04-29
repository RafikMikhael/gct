package main

import (
	"context"
	"errors"
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
	muxRouter *mux.Router
	mu        sync.Mutex // mutex used to access jobs map
	jobs      map[string]*Job
	bStopped  bool
	Port      string

	// stuff needed for graceful shutdown of server
	cancel      context.CancelFunc
	ctx         context.Context
	StopSignals chan os.Signal

	// look-up tables
	bitRate   [NumberOfQualityLvls][NumberOfRenditions]int //[quality][renditionIdx]
	horizW    [NumberOfRenditions]int                      //rendition target width in pixels
	vertH     [NumberOfRenditions]int                      //rendition target height in pixels
	sleepTime [NumberOfRenditions]int                      //in seconds
}

// Initialize - initialize App fields and allocate all needed memory
func (app *App) Initialize(portNum *int) error {
	if *portNum <= 1024 || *portNum >= 49151 {
		return errors.New("port number should be between 1025 and 49150")
	}

	app.bitRate = [NumberOfQualityLvls][NumberOfRenditions]int{
		{160, 360, 1930, 4080, 7000},
		{145, 300, 1600, 3400, 5800},
		{120, 280, 1400, 3080, 4500},
	}
	app.horizW = [NumberOfRenditions]int{640, 768, 960, 1280, 1920}
	app.vertH = [NumberOfRenditions]int{360, 432, 540, 720, 1080}
	app.sleepTime = [NumberOfRenditions]int{10, 20, 30, 40, 50}
	app.jobs = make(map[string]*Job)
	app.Port = ":" + strconv.Itoa(*portNum)

	app.StopSignals = make(chan os.Signal, 1)
	signal.Notify(app.StopSignals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	return nil
}

// Run - run the application (main go routine running forever)
func (app *App) Run() {
	// graceful shutdown can be triggered by /terminate or by ctrl+c
	// https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97
	app.ctx, app.cancel = context.WithCancel(context.Background())
	go func() {
		<-app.StopSignals
		app.cancel()
	}()

	app.muxRouter = mux.NewRouter().StrictSlash(true)
	app.muxRouter.HandleFunc("/_health", http.HandlerFunc(app.isupHandler))
	app.muxRouter.HandleFunc("/api/v1/terminate", app.terminate)
	app.muxRouter.HandleFunc("/api/v1/job/{quality}", app.triggerJobs)
	app.muxRouter.HandleFunc("/api/v1/probe/{hash}", app.probeHash)
	app.muxRouter.HandleFunc("/api/v1/monitor", app.monitor)

	server := http.Server{
		Addr:    app.Port,
		Handler: app.muxRouter,
	}
	// we can add any middleware here like corsHandler or LoggingHandler
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	fmt.Printf("transcode server starting up\n")
	<-app.ctx.Done()
	fmt.Printf("server stopped, via api = %v\n", app.bStopped)

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed:%+s", err)
	}

	fmt.Printf("transcode server shut down properly\n")
}

// terminate - cleanly close all go routines and recover resources
func (app *App) terminate(w http.ResponseWriter, r *http.Request) {
	// verify the verb used
	if r.Method != "GET" {
		// 405
		app.JsonHttpResponse(w, http.StatusMethodNotAllowed, "error", r.Method)
		return
	}
	app.bStopped = true
	app.JsonHttpResponse(w, http.StatusOK, "termination", "started")
}

// JsonHttpResponse - format a JSON http response using http Statu Code and key-value body
func (app *App) JsonHttpResponse(w http.ResponseWriter, code int, key, value string) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("{\"%s\":%s}", key, value)))
}

// isupHandler - serve responses to _health GET requests needed by AWS
func (app *App) isupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Should actually check something to verify the server is not stuck
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"up":true}`))
}
