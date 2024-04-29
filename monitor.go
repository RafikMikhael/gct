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
	keys := maps.Keys(app.Jobs)
	defer app.mu.Unlock()
	fmt.Fprintf(w, "ongoing hashes=%v\n", keys)
	if app.BStopped && len(app.Jobs) == 0 {
		app.StopSignals <- syscall.SIGINT
	}
}
