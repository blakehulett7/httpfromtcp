package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blakehulett7/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, syscall.SIGINT, syscall.SIGTERM)
	<-sig_chan
	log.Println("Server gracefully stopped")
}
