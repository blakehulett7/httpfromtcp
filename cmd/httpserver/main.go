package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/blakehulett7/httpfromtcp/internal/headers"
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

	if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, r)
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

func proxyHandler(w *response.Writer, r *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	w.Headers = headers.Headers{
		"content-type":      "text/plain",
		"transfer-encoding": "chunked",
		"Trailer":           "X-Content-SHA256, X-Content-Length",
	}
	w.WriteHeaders()

	route := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")
	path := fmt.Sprintf("https://httpbin.org/%s", route)
	res, err := http.Get(path)
	if err != nil {
		apiErr := response.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
		w.WriteError(apiErr, "text/plain")
		return
	}

	full_response := make([]byte, 0)
	for {
		buffer := make([]byte, 1024)

		var length int
		length, err = res.Body.Read(buffer)
		if err != nil {
			break
		}
		fmt.Println(length)

		if length > 0 {
			w.WriteChunkedBody(buffer[:length])
			full_response = append(full_response, buffer[:length]...)
		}
	}

	if err != io.EOF {
		fmt.Println(err)
		return
	}

	w.WriteChunkedBodyDone()

	hash := sha256.Sum256(full_response)
	content_length := len(full_response)

	trailers := headers.Headers{
		"X-Content-SHA256": fmt.Sprintf("%x", hash),
		"X-Content-Length": fmt.Sprint(content_length),
	}
	w.WriteTrailers(trailers)
}
