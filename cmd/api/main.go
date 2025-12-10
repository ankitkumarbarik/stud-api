package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "github.com/ankitkumarbarik/rest-api/internal/http"
)

func main() {
	router := httpapi.NewRouter()

	server := http.Server{
		Addr:    ":5000",
		Handler: router,
	}

	// channel to listen OS signals
	stop := make(chan os.Signal, 1)

	// register signals
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// server run goroutine
	go func() {
		log.Println("server listing on PORT", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// wait for signal
	<-stop
	log.Println("shutting down server...")

	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown failed:", err)
	}

	log.Println("server stopped gracefully")
}

// go build -o api.exe ./cmd/api; .\api.exe
