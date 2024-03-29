package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	port       string
	debug      bool
	messages   msgFlags
	allowSleep bool
)

func init() {
	flag.Var(&messages, "msg", "custom message to print out. can be used multiple times")
	flag.StringVar(&port, "p", "8080", "Port to listen on")
	flag.BoolVar(&debug, "d", true, "debug. Print all requests")
	flag.BoolVar(&allowSleep, "allow-sleep", true, "Allows ?sleep=<duration> query parameter")
}

func main() {
	flag.Parse()

	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})

	http.HandleFunc("/api/gowhoami/log", func(_ http.ResponseWriter, _ *http.Request) {
		logrus.WithFields(logrus.Fields{
			"field1": "test",
		}).Info("Test logging")
	})
	http.HandleFunc("/", handler)
	log.Println("Starting server on port: " + port)

	server, shutdown := NewServerWithTimeout(10 * time.Second)
	server.Handler = http.DefaultServeMux
	server.Addr = ":" + port

	log.Println(server.ListenAndServe())

	<-shutdown
	log.Println("Stopped")
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	for _, msg := range messages {
		fmt.Fprintf(w, "%s\n", msg)
	}
	fmt.Fprintf(w, "Hostname: %s\n", hostname)
	fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
	fmt.Fprintf(w, "Path: %s\n", r.URL.Path)

	keys := make([]string, 0, len(r.Header))
	for key := range r.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "Headers:\n")
	for _, v := range keys {
		fmt.Fprintf(w, "%s:%s\n", v, r.Header.Get(v))
	}

	if debug {
		log.Printf("| %s | %s | %s", r.RemoteAddr, r.Method, r.RequestURI)
		log.Printf("Host: %s\n", r.Host)
		for _, v := range keys {
			log.Printf("%s:%s\n", v, r.Header.Get(v))
		}
	}

	if allowSleep {
		strTime := r.URL.Query().Get("sleep")
		if strTime == "" {
			return
		}
		dur, err := time.ParseDuration(strTime)
		if err != nil {
			logrus.Error(err)
			fmt.Fprintf(w, "error: %s\n", err.Error())
			return
		}
		time.Sleep(dur)
	}
}

func NewServerWithTimeout(t time.Duration) (*http.Server, chan struct{}) {
	shutdown := make(chan struct{})
	srv := &http.Server{}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("gograce: Shutdown Server ...")

		time.Sleep(5 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), t)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("gograce: error server shutdown:", err)
		}
		close(shutdown)
		log.Println("gograce: server exited")
	}()

	return srv, shutdown
}
