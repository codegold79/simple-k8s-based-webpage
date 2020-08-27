package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Requests struct {
	SleepTime, StatusCode int
}

func main() {
	http.HandleFunc("/300/", statusRedirectHandler)
	http.HandleFunc("/200/", statusOkHandler)
	http.HandleFunc("/400/", statusUserErrHandler)
	http.HandleFunc("/500/", statusServerErrHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)

	fmt.Println("Serving...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func statusOkHandler(w http.ResponseWriter, r *http.Request) {
	var sleepTime int

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{sleepTime, 200}

	executeTemplate(w, req)
}

func statusRedirectHandler(w http.ResponseWriter, r *http.Request) {
	var sleepTime int

	w.WriteHeader(http.StatusMultipleChoices)

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{sleepTime, 300}

	executeTemplate(w, req)
}

func statusUserErrHandler(w http.ResponseWriter, r *http.Request) {
	var sleepTime int

	w.WriteHeader(http.StatusTeapot)

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{sleepTime, 418}

	executeTemplate(w, req)
}

func statusServerErrHandler(w http.ResponseWriter, r *http.Request) {
	var sleepTime int

	w.WriteHeader(http.StatusInsufficientStorage)

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{sleepTime, 507}

	executeTemplate(w, req)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func executeTemplate(w http.ResponseWriter, req Requests) {
	tmpl := template.Must(
		template.ParseFiles("go-template.html"),
	)
	err := tmpl.ExecuteTemplate(w, "go-template.html", req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func waitForSleep(w http.ResponseWriter, vals []string) int {
	sleepTime, err := strconv.Atoi(vals[0])
	if err != nil {
		http.Error(
			w,
			"sleep request value needs to be an integer: "+err.Error(),
			http.StatusBadRequest,
		)
	}

	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	return sleepTime
}
