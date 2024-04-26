package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Quality int
type RenditionIdx int

const (
	HIGH Quality = iota
	MEDIUM
	LOW
)

const (
	RATE0 RenditionIdx = iota
	RATE1
	RATE2
	RATE3
	RATE4
)

type Job struct {
	Qual  Quality
	Hash  string
	wg    sync.WaitGroup
	sizes chan int
}

// CreateJob - trigger the jobs encoding the input path to output path according to quality
func (app *App) CreateJob(w http.ResponseWriter, r *http.Request) {
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

	job := &Job{
		sizes: make(chan int, 5),
	}
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

	fmt.Printf("creating Jobs for bw=%v, in=%v, out=%v, w=%d, h=%d\n", bandwidth, inputPath, outputPath, width, height)
	//fmt.Fprintf(w, "bw=%v, in=%v, out=%v, w=%d, h=%d\n", bandwidth, inputPath, outputPath, width, height)

	// produce the hash for tracking
	job.Hash = hashThis(inputPath, outputPath)

	// trigger the 5 goroutines
	job.wg.Add(5)
	go app.transcodeRendition(job, inputPath, outputPath+"0", RATE0)
	go app.transcodeRendition(job, inputPath, outputPath+"1", RATE1)
	go app.transcodeRendition(job, inputPath, outputPath+"2", RATE2)
	go app.transcodeRendition(job, inputPath, outputPath+"3", RATE3)
	go app.transcodeRendition(job, inputPath, outputPath+"4", RATE4)

	// closer
	go func() {
		job.wg.Wait()
		close(job.sizes)
	}()

	var total int
	for size := range job.sizes {
		total += size
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("{\"id\":%s,\"size\":%d}", job.Hash[:], total)))
}

func hashThis(input, output string) string {
	timeInputString := time.Now().String() + input + output
	hash := md5.Sum([]byte(timeInputString))
	return hex.EncodeToString(hash[:])
}
func (app *App) transcodeRendition(job *Job, inputPath, outputPath string, brIdx RenditionIdx) {
	defer job.wg.Done()
	fmt.Printf("starting %s, outname=%s, w=%4d, h=%4d, bitrate=%4d, duration=%d\n",
		job.Hash, outputPath, app.horizW[brIdx], app.vertH[brIdx], app.bitRate[job.Qual][brIdx], app.sleepTime[brIdx])
	time.Sleep(time.Duration(app.sleepTime[brIdx] * int(time.Second)))
	fmt.Printf("finished %s, outname=%s, w=%4d, h=%4d, bitrate=%4d, duration=%d\n",
		job.Hash, outputPath, app.horizW[brIdx], app.vertH[brIdx], app.bitRate[job.Qual][brIdx], app.sleepTime[brIdx])
	job.sizes <- app.sleepTime[brIdx] // simulation for size in bytes
}
