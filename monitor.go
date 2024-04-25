package main

import (
	"log"
	"net/http"
)

// Monitor - go routine that supports non-blocking stats for the App resources
func (app *App) Monitor() {
	http.HandleFunc("/", app.MonitorHandler)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
func (app *App) MonitorHandler(w http.ResponseWriter, r *http.Request) {
}
