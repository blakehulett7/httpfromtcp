package server

import (
	"io"

	"github.com/blakehulett7/httpfromtcp/internal/request"
	"github.com/blakehulett7/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Error      error
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func writeError(w io.Writer, err *HandlerError) {
	response.WriteStatusLine(w, response.StatusCode(err.StatusCode))
	message := []byte(err.Error.Error())
	response.WriteHeaders(w, response.GetDefaultHeaders(len(message)))
	w.Write(message)
}
