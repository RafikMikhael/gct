package main

import (
	"fmt"
	"net/http"
	"syscall"

	"golang.org/x/exp/maps"
)

// Monitor - go routine that supports non-blocking stats for the App resources
func (app *App) Monitor(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	keys := maps.Keys(app.jobs)
	defer app.mu.Unlock()
	fmt.Fprintf(w, "ongoing hashes=%v\n", keys)
	if app.bStopped && len(app.jobs) == 0 {
		app.StopSignals <- syscall.SIGINT
	}
}
