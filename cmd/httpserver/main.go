package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blakehulett7/httpfromtcp/internal/request"
	"github.com/blakehulett7/httpfromtcp/internal/response"
	"github.com/blakehulett7/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
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

func handler(w io.Writer, r *request.Request) *server.HandlerError {
	if r.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Error:      fmt.Errorf("Your problem is not my problem\n"),
		}
	}

	if r.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Error:      fmt.Errorf("Woopsie, my bad\n"),
		}
	}

	w.Write([]byte("All good, frfr\n"))
	return nil
}
