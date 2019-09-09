package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/willgorman/homebridge-unicorn-hat/internal/pkg/fake"
	"github.com/willgorman/homebridge-unicorn-hat/internal/pkg/light"

	log "github.com/sirupsen/logrus"
)

var theLight light.Light

func main() {

	theLight = &fake.FakeLight{}

	r := mux.NewRouter()
	r.HandleFunc("/api/switch", SwitchStatusHandler)
	r.HandleFunc("/api/switch/on", SwitchHandler(true))
	r.HandleFunc("/api/switch/off", SwitchHandler(false))
	r.HandleFunc("/api/brightness", BrightnessHandler)
	r.HandleFunc("/api/brightness/{value}", SetBrightnessHandler)
	r.HandleFunc("/api/color", ColorHandler)
	r.HandleFunc("/api/color/{value}", SetColorHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}

func SwitchStatusHandler(w http.ResponseWriter, r *http.Request) {

	log.Infof("Getting switch status")
	on, err := theLight.IsOn()
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
	}
	var status string
	if on {
		status = "0"
	} else {
		status = "1"
	}
	io.WriteString(w, status)
}

func SwitchHandler(on bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		if on {
			err = theLight.TurnOn()
		} else {
			err = theLight.TurnOff()
		}
		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, err.Error())
			return
		}

		w.WriteHeader(200)

	}
}

func BrightnessHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Getting brightness %s", 100)
	w.WriteHeader(200)
	io.WriteString(w, "100")
}

func SetBrightnessHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	brightness := vars["value"]
	log.Infof("Setting brightness to %s", brightness)
	w.WriteHeader(200)
}

func ColorHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Getting color")
	w.WriteHeader(200)
	io.WriteString(w, "000000")
}

func SetColorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	color := vars["value"]
	log.Infof("Setting color", color)
	w.WriteHeader(200)
}
