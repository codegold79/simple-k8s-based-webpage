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

const (
	status2xxImagePath = "graphics/MariaLetta_gophers-basking-in-sun.png"
	status3xxImagePath = "graphics/MariaLetta_gopher-at-crossroads.png"
	status4xxImagePath = "graphics/MariaLetta_gopher-on-mobile.png"
	status5xxImagePath = "graphics/MariaLetta_gopher-in-flames.png"
)

// Requests stores the amount of time to sleep before a response
// as well as the status code to return.
type Requests struct {
	SleepTime, StatusCode           int
	GopherImagePath, GopherImageAlt string
}

func main() {
	http.HandleFunc("/", envStatusHandler)
	http.Handle("/graphics/", http.StripPrefix("/graphics/", http.FileServer(http.Dir("graphics"))))

	http.HandleFunc("/200/", status200Handler)
	http.Handle("/200/graphics/", http.StripPrefix("/200/graphics/", http.FileServer(http.Dir("graphics"))))

	http.HandleFunc("/300/", status300Handler)
	http.Handle("/300/graphics/", http.StripPrefix("/300/graphics/", http.FileServer(http.Dir("graphics"))))

	http.HandleFunc("/400/", status400Handler)
	http.Handle("/400/graphics/", http.StripPrefix("/400/graphics/", http.FileServer(http.Dir("graphics"))))

	http.HandleFunc("/500/", status500Handler)
	http.Handle("/500/graphics/", http.StripPrefix("/500/graphics/", http.FileServer(http.Dir("graphics"))))

	http.HandleFunc("/favicon.ico", faviconHandler)

	fmt.Println("Serving on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func envStatusHandler(w http.ResponseWriter, r *http.Request) {
	code := 200

	// Overwrite code with environmental variable if it has a valid value.
	if envVarCode() > 0 {
		code = envVarCode()
	}

	handleCode(w, r, code)
}

// validCode grabs an environmental variable and stores it as an integer.
func envVarCode() int {
	envCode := os.Getenv("STATUSCODE")
	code, err := strconv.Atoi(envCode)

	if err != nil {
		return -1
	}

	return code
}

func status200Handler(w http.ResponseWriter, r *http.Request) {
	code := 200

	handleCode(w, r, code)
}

func status300Handler(w http.ResponseWriter, r *http.Request) {
	code := 300

	handleCode(w, r, code)
}

func status400Handler(w http.ResponseWriter, r *http.Request) {
	code := 400

	handleCode(w, r, code)
}

func status500Handler(w http.ResponseWriter, r *http.Request) {
	code := 500

	handleCode(w, r, code)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "graphics/favicon.ico")
}

func handleCode(w http.ResponseWriter, r *http.Request, code int) {
	var sleepTime int

	sleepTimes, ok := r.URL.Query()["sleep"]
	if ok {
		sleepTime = waitForSleep(w, sleepTimes)
	}

	req := Requests{
		SleepTime:  sleepTime,
		StatusCode: code,
	}

	w.WriteHeader(code)
	setStatusImage(&req, code)

	executeTemplate(w, req)
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

func setStatusImage(req *Requests, code int) {
	switch {
	case code >= 500:
		req.GopherImagePath = status5xxImagePath
		req.GopherImageAlt = "Alt: Image representing 5xx status code."
	case code >= 400:
		req.GopherImagePath = status4xxImagePath
		req.GopherImageAlt = "Alt: Image representing 4xx status code."
	case code >= 300:
		req.GopherImagePath = status3xxImagePath
		req.GopherImageAlt = "Alt: Image representing 3xx status code."
	case code >= 200:
		req.GopherImagePath = status2xxImagePath
		req.GopherImageAlt = "Alt: Image representing 2xx status code."
	}
}
