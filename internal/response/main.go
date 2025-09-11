package response

import (
	"fmt"
	"io"

	h "github.com/blakehulett7/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(content_length int) h.Headers {
	return h.Headers{
		"Content-Length": fmt.Sprint(content_length),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func WriteHeaders(w io.Writer, headers h.Headers) error {
	for key, value := range headers {
		_, err := w.Write(fmt.Appendf(nil, "%s: %s\r\n", key, value))
		if err != nil {
			return err
		}
	}

	w.Write([]byte("\r\n"))

	return nil
}

func WriteStatusLine(w io.Writer, status_code StatusCode) error {
	switch status_code {

	case StatusOK:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusBadRequest:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusInternalServerError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		w.Write(fmt.Appendf(nil, "HTTP/1.1 %d\r\n", status_code))
		return fmt.Errorf("Unsupported status code")
	}

	return nil
}
