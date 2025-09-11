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
		_, err := w.Write(fmt.Appendf(nil, "%s: %s", key, value))
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteStatusLine(w io.Writer, status_code StatusCode) error {
	switch status_code {

	case StatusOK:
		fmt.Fprintf(w, "HTTP/1.1 200 OK")
	case StatusBadRequest:
		fmt.Fprintf(w, "HTTP/1.1 400 Bad Request")
	case StatusInternalServerError:
		fmt.Fprintf(w, "HTTP/1.1 500 Internal Server Error")
	default:
		fmt.Fprintf(w, "HTTP/1.1 %d", status_code)
		return fmt.Errorf("Unsupported status code")
	}

	return nil
}
