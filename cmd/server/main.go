package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/willgorman/homebridge-pi-light/internal/pkg/blinkt"
	"github.com/willgorman/homebridge-pi-light/internal/pkg/fake"
	"github.com/willgorman/homebridge-pi-light/internal/pkg/light"
	"github.com/willgorman/homebridge-pi-light/internal/pkg/unicorn"
)

var theLight light.Light

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("hpilight")
	viper.SetDefault("log_level", "warn")
	viper.SetDefault("light_type", "unicorn")

	log.SetReportCaller(true)
	switch viper.GetString("log_level") {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	}

}

func newLight() light.Light {
	lightType := viper.GetString("light_type")
	switch lightType {
	case "unicorn":
		var err error
		l, err := unicorn.NewUnicornLight()
		if err != nil {
			panic(err)
		}
		return l
	case "blinkt":
		return blinkt.New()
	default:
		log.Info("Creating fake light")
		return &fake.FakeLight{}
	}

}

func main() {
	log.SetReportCaller(true)

	theLight = newLight()
	log.Infof("The light: %v", theLight)

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("OK!"))
	})
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
	theLight.TurnOff()
	os.Exit(0)

}

func SwitchStatusHandler(w http.ResponseWriter, r *http.Request) {

	log.Infof("Getting switch status for light %v", theLight)
	on, err := theLight.IsOn()
	if err != nil {
		w.WriteHeader(500)
		log.Errorf("Failed to get status: %v", err)
		io.WriteString(w, err.Error())
	}
	var status string
	if on {
		log.Infof("The light is on")
		status = "1"
	} else {
		log.Infof("The light is off")
		status = "0"
	}
	log.Infof("Returning status: %v", status)
	io.WriteString(w, status)
}

func SwitchHandler(on bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		if on {
			err = theLight.TurnOn()
			log.Info("Turned light on")
		} else {
			err = theLight.TurnOff()
			log.Info("Turned light off")
		}
		if err != nil {
			log.Errorf("Failed to turn light (%v): %v", on, err)
			w.WriteHeader(500)
			io.WriteString(w, err.Error())
			return
		}

		w.WriteHeader(200)

	}
}

func BrightnessHandler(w http.ResponseWriter, r *http.Request) {

	b, err := theLight.GetBrightness()
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
		return
	}

	log.Infof("Getting brightness %d", b)
	io.WriteString(w, strconv.Itoa(int(b)))
}

func SetBrightnessHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	brightness := vars["value"]
	log.Infof("Setting brightness to %s", brightness)
	bint, err := strconv.ParseUint(brightness, 10, 8)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
	}

	err = theLight.SetBrightness(uint(bint))
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
	}

	w.WriteHeader(200)
}

func ColorHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Getting color")
	c, err := theLight.GetColor()
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
	}
	log.Infof("Getting color: %s (%s)", c.ToHexString(), c)
	w.WriteHeader(200)
	io.WriteString(w, c.ToHexString())
}

func SetColorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	color := vars["value"]
	log.Infof("Setting color %v", color)

	rHex := color[0:2]
	gHex := color[2:4]
	bHex := color[4:6]
	red, err := strconv.ParseInt(rHex, 16, 64)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "Invalid color")
		log.Errorf("Failed to set color: %v", err)
		return
	}
	green, err := strconv.ParseInt(gHex, 16, 64)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "Invalid color")
		log.Errorf("Failed to set color: %v", err)
		return
	}
	blue, err := strconv.ParseInt(bHex, 16, 64)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "Invalid color")
		log.Errorf("Failed to set color: %v", err)
		return
	}

	log.Info("Setting color to light")
	err = theLight.SetColor(uint(red), uint(green), uint(blue))
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprintf("Failed to set color: %v", err))
		log.Errorf("Failed to set color: %v", err)
		return
	}
	w.WriteHeader(200)
}
