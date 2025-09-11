package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blakehulett7/httpfromtcp/internal/request"
	"github.com/blakehulett7/httpfromtcp/internal/response"
	"github.com/blakehulett7/httpfromtcp/internal/server"
)

const port = 42069
const dir = "./cmd/httpserver"

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

func handler(w *response.Writer, r *request.Request) {
	if r.RequestLine.RequestTarget == "/yourproblem" {
		path := fmt.Sprintf("%s/html/bad_request.html", dir)
		bad_request, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		h_err := response.HandlerError{
			StatusCode: response.StatusBadRequest,
			Error:      fmt.Errorf(string(bad_request)),
		}

		w.WriteError(h_err, "text/html")
		return
	}

	if r.RequestLine.RequestTarget == "/myproblem" {
		path := fmt.Sprintf("%s/html/internal_error.html", dir)
		internal_error, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		h_err := response.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Error:      fmt.Errorf(string(internal_error)),
		}

		w.WriteError(h_err, "text/html")
		return
	}

	path := fmt.Sprintf("%s/html/ok.html", dir)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.WriteStatusLine(response.StatusOK)
	w.Headers = response.GetDefaultHeaders(len(data))
	w.Headers.Set("content-type", "text/html")
	w.WriteHeaders()
	w.WriteBody(data)
}
