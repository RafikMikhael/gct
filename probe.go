package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ProbeHash - Get info about a specific job with a knonw hash
func (app *App) ProbeHash(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	job := app.Jobs[hash]
	if job != nil {
		job.mu.Lock()
		doneStr := strings.Join(job.DoneRenditions[:], ",")
		job.mu.Unlock()
		app.JsonHttpResponse(w, http.StatusOK, "done", doneStr)
	} else {
		// 204
		app.JsonHttpResponse(w, http.StatusNotFound, "done", hash)
	}

}
