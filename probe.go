package main

import (
	"fmt"
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
		w.Header().Set("content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"done\":%s}", doneStr)))
	} else {
		// 204
		app.ErrorResponse(w, http.StatusNotFound, "done", hash)
	}

}
