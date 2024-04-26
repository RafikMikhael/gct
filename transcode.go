package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Job struct {
	Qual Quality
	Hash string
}

// Transcode - transcode the input path to output path according to quality
func (app *App) Transcode(w http.ResponseWriter, r *http.Request) {
	bandwidth := mux.Vars(r)["quality"]
	inputPath := r.URL.Query().Get("inputpath")
	outputPath := r.URL.Query().Get("outputpath")
	width, errW := strconv.Atoi(r.URL.Query().Get("w"))
	height, errH := strconv.Atoi(r.URL.Query().Get("h"))

	// verify the verb used
	if r.Method != "POST" {
		app.ErrorResponse(w, http.StatusMethodNotAllowed, r.Method)
		return
	}

	job := &Job{}
	// Do some validation
	switch strings.ToLower(bandwidth) {
	case "high":
		job.Qual = HIGH
	case "medium":
		job.Qual = MEDIUM
	case "low":
		job.Qual = LOW
	default:
		app.ErrorResponse(w, http.StatusBadRequest, "quality")
		return
	}

	// we only support wigth between 640 and 4K
	if errW != nil || width < 640 || width > 3840 {
		app.ErrorResponse(w, http.StatusBadRequest, "width")
		return
	}
	// we only support height between 480 and 4K
	if errH != nil || height < 480 || height > 2176 {
		app.ErrorResponse(w, http.StatusBadRequest, "height")
		return
	}

	//fmt.Printf("bw=%v, in=%v, out=%v, w=%d, h=%d\n", bandwidth, inputPath, outputPath, width, height)
	//fmt.Fprintf(w, "bw=%v, in=%v, out=%v, w=%d, h=%d\n", bandwidth, inputPath, outputPath, width, height)

	// proceed with the transcode, produce the hash for tracking
	job.Hash = hashThis(inputPath, outputPath)
	app.ValidResponse(w, string(job.Hash[:]))
}

func hashThis(input, output string) string {
	timeInputString := time.Now().String() + input + output
	hash := md5.Sum([]byte(timeInputString))
	return hex.EncodeToString(hash[:])
}
