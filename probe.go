package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ProbeHash - Get info about a specific job with a knonw hash
func (app *App) ProbeHash(w http.ResponseWriter, r *http.Request) {
	// verify the verb used
	if r.Method != "GET" {
		// 405
		app.JsonHttpResponse(w, http.StatusMethodNotAllowed, "error", r.Method)
		return
	}
	hash := mux.Vars(r)["hash"]
	job := app.jobs[hash]
	if job != nil {
		job.mu.Lock()
		doneStr := strings.Join(job.DoneRenditions[:], ",")
		job.mu.Unlock()
		app.JsonHttpResponse(w, http.StatusOK, "done", doneStr)
	} else {
		// 204
		app.JsonHttpResponse(w, http.StatusNoContent, "done", hash)
	}

}
