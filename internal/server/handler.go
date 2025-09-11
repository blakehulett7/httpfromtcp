package server

import (
	"github.com/blakehulett7/httpfromtcp/internal/request"
	"github.com/blakehulett7/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)
