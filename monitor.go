package main

import (
	"fmt"
	"log"
	"net/http"
	"syscall"

	"golang.org/x/exp/maps"
)

// Monitor - go routine that supports non-blocking stats for the App resources
func (app *App) Monitor() {
	http.HandleFunc("/", app.MonitorHandler)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
func (app *App) MonitorHandler(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	keys := maps.Keys(app.Jobs)
	app.mu.Unlock()
	fmt.Fprintf(w, "ongoing hashes=%v\n", keys)
	if app.BStopped && len(app.Jobs) == 0 {
		app.StopSignals <- syscall.SIGINT
	}
}
