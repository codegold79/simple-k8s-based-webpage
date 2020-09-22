package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Requests stores the amount of time to sleep before a response
// as well as the status code to return.
type Requests struct {
	SleepTime, StatusCode int
}

func main() {
	http.HandleFunc("/", envStatusHandler)
	http.HandleFunc("/200/", status200Handler)
	http.HandleFunc("/300/", status300Handler)
	http.HandleFunc("/400/", status400Handler)
	http.HandleFunc("/500/", status500Handler)
	http.HandleFunc("/favicon.ico", faviconHandler)

	fmt.Println("Serving on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func envStatusHandler(w http.ResponseWriter, r *http.Request) {
	code := validCode()

	req := Requests{
		SleepTime:  0,
		StatusCode: code,
	}

	if code > 0 {
		w.WriteHeader(code)
	}

	executeTemplate(w, req)
}

// validCode grabs an environmental variable and stores it as an integer.
func validCode() int {
	envCode := os.Getenv("STATUSCODE")
	code, err := strconv.Atoi(envCode)

	if err != nil {
		log.Printf("STATUSCODE env var needs to be an integer between 200-599, inclusive: %v\n", err)
		return -1
	}

	return code
}

func status200Handler(w http.ResponseWriter, r *http.Request) {
	var sleepTime int

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{sleepTime, 200}

	executeTemplate(w, req)
}

func status300Handler(w http.ResponseWriter, r *http.Request) {
	var sleepTime int

	w.WriteHeader(http.StatusMultipleChoices)

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{sleepTime, 300}

	executeTemplate(w, req)
}

func status400Handler(w http.ResponseWriter, r *http.Request) {
	var sleepTime int

	w.WriteHeader(http.StatusBadRequest)

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{sleepTime, 400}

	executeTemplate(w, req)
}

func status500Handler(w http.ResponseWriter, r *http.Request) {
	var sleepTime int

	w.WriteHeader(http.StatusInternalServerError)

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{sleepTime, 500}

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
