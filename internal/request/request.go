package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	raw, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	parts := strings.Split(string(raw), "\r\n")
	raw_request_line := parts[0]
	request_line, err := parseRequestLine(raw_request_line)
	if err != nil {
		return &Request{}, err
	}

	return &Request{request_line}, nil
}

func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("malformed request line")
	}

	method := parts[0]

	return RequestLine{
		HttpVersion:   parts[2],
		RequestTarget: parts[1],
		Method:        parts[0],
	}, nil
}
