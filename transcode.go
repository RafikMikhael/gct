package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
	Qual           Quality
	Hash           string
	mu             sync.Mutex
	DoneRenditions []string
	wg             sync.WaitGroup
	sizes          chan int
}

// TriggerJobs - trigger the jobs encoding the input path to output path according to quality
func (app *App) TriggerJobs(w http.ResponseWriter, r *http.Request) {
	// server instructed to stop as soon as its managed jobs are done
	if app.BStopped {
		// 503
		app.JsonHttpResponse(w, http.StatusServiceUnavailable, "termination", "started")
		return
	}

	bandwidth := mux.Vars(r)["quality"]
	inputPath := r.URL.Query().Get("inputpath")
	outputPath := r.URL.Query().Get("outputpath")
	width, errW := strconv.Atoi(r.URL.Query().Get("w"))
	height, errH := strconv.Atoi(r.URL.Query().Get("h"))

	// verify the verb used
	if r.Method != "POST" {
		// 405
		app.JsonHttpResponse(w, http.StatusMethodNotAllowed, "error", r.Method)
		return
	}

	job := &Job{
		DoneRenditions: []string{},
		sizes:          make(chan int, NumberOfRenditions),
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
		app.JsonHttpResponse(w, http.StatusBadRequest, "error", "quality")
		return
	}

	// we only support wigth between 640 and 4K
	if errW != nil || width < 640 || width > 3840 {
		app.JsonHttpResponse(w, http.StatusBadRequest, "error", "width")
		return
	}
	// we only support height between 480 and 4K
	if errH != nil || height < 480 || height > 2176 {
		app.JsonHttpResponse(w, http.StatusBadRequest, "error", "height")
		return
	}

	fmt.Printf("creating Jobs for quality=%v, in=%v, out=%v, w=%d, h=%d\n", bandwidth, inputPath, outputPath, width, height)

	// produce the hash for tracking
	job.Hash = hashThis(inputPath, outputPath)

	// as we reach this point, the job has valid inputs and a valid hash, we add it to Jobs map and start processing
	app.mu.Lock()
	app.Jobs[job.Hash] = job
	app.mu.Unlock()

	go app.startWorkers(job, inputPath, outputPath)

	// write the response right away, so client can use the hash for probing
	app.JsonHttpResponse(w, http.StatusOK, "id", job.Hash[:])
}

func hashThis(input, output string) string {
	timeInputString := time.Now().String() + input + output
	hash := md5.Sum([]byte(timeInputString))
	return hex.EncodeToString(hash[:])
}

func (app *App) transcodeRendition(job *Job, inputPath, outputPath string, brIdx RenditionIdx) {
	defer job.wg.Done()

	myW := app.horizW[brIdx]
	myH := app.vertH[brIdx]
	fmt.Printf("starting %s, outname=%s, w=%4d, h=%4d, bitrate=%4d, duration=%d\n",
		job.Hash, outputPath, myW, myH, app.bitRate[job.Qual][brIdx], app.sleepTime[brIdx])
	time.Sleep(time.Duration(app.sleepTime[brIdx] * int(time.Second)))
	fmt.Printf("finished %s, outname=%s, w=%4d, h=%4d, bitrate=%4d, duration=%d\n",
		job.Hash, outputPath, myW, myH, app.bitRate[job.Qual][brIdx], app.sleepTime[brIdx])

	job.mu.Lock()
	job.DoneRenditions = append(job.DoneRenditions, strconv.Itoa(myW)+`x`+strconv.Itoa(myH))
	job.mu.Unlock()

	job.sizes <- app.sleepTime[brIdx] // simulate size in KBytes to be the same as duration
}

func (app *App) startWorkers(job *Job, inputPath, outputPath string) {
	// trigger the NumberOfRenditions goroutines
	job.wg.Add(NumberOfRenditions)
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
	fmt.Printf("---job id %s finished with total size %d\n", job.Hash, total)
	app.mu.Lock()
	defer app.mu.Unlock()
	delete(app.Jobs, job.Hash)
	if app.BStopped && len(app.Jobs) == 0 {
		app.StopSignals <- syscall.SIGINT
	}
}
