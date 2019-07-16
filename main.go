package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
)

var location *time.Location

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setupRouter(fetcher APIFetcher) {
	http.HandleFunc("/v1/live", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		live := fetcher.Live()
		if err := json.NewEncoder(w).Encode(live); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	})
}

func main() {
	timezone := os.Getenv("API_TIMEZONE")
	if timezone == "" {
		log.Fatal("missing API_TIMEZONE env var")
	}
	l, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not load timezone"))
	}
	location = l

	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		log.Fatal("missing API_HOST env var")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")

	fetcher := &apiFetcher{
		host: apiHost,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}

	setupRouter(fetcher)

	// Handle Ctrl+C
	var quit chan os.Signal
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	tick := time.NewTicker(time.Second).C
	go func() {
		for {
			select {
			case <-quit:
				log.Println("exiting...")
				os.Exit(0)
			case <-tick:
				if err := fetcher.PollAPI(); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	log.Printf("listening on %s:%s\n", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil))
}
