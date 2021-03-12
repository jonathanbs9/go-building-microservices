package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jonathanbs9/go-building-microservices/handlers"
)

func main() {
	// Logger
	l := log.New(os.Stdout, "product-api | ", log.LstdFlags)

	// Handlers
	//hh := handlers.NewHello(l)
	//gh := handlers.NewGoodbye(l)
	ph := handlers.NewProducts(l)

	// ServeMux
	sm := http.NewServeMux()
	//sm.Handle("/", hh)
	//sm.Handle("/goodbye", gh)
	sm.Handle("/products", ph)

	// Server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		l.Println("Starting Server on port :9090 ...")

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until signal is received
	sig := <-c
	log.Println("Got signal! => ", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
