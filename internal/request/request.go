package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
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
		return RequestLine{}, fmt.Errorf("malformed request line: %v", line)
	}

	method := parts[0]
	if !containsOnlyCapitalLetters(method) {
		return RequestLine{}, fmt.Errorf("method must contain capital letters only: %v", method)
	}

	http_version_literal := parts[2]
	version_parts := strings.Split(http_version_literal, "/")
	http_version := version_parts[1]
	if len(version_parts) != 2 || http_version != "1.1" {
		return RequestLine{}, fmt.Errorf("malformed or incorrect http version: %v", http_version_literal)
	}

	return RequestLine{
		HttpVersion:   http_version,
		RequestTarget: parts[1],
		Method:        method,
	}, nil
}

func containsOnlyCapitalLetters(s string) bool {
	for _, c := range s {
		if !unicode.IsUpper(c) {
			return false
		}
	}

	return true
}
