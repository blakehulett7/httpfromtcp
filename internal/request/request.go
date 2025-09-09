package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	h "github.com/blakehulett7/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     h.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := NewBuffer()

	raw_request_line, err := buffer.readLine(reader)
	request_line, err := parseRequestLine(raw_request_line)
	if err == io.EOF {
		return &Request{}, fmt.Errorf("malformed request... no host")
	}

	if err != nil {
		return &Request{}, err
	}

	headers, err := parseHeaders(&buffer, reader)
	if err != nil {
		return &Request{}, err
	}

	return &Request{
		RequestLine: request_line,
		Headers:     headers,
	}, nil
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

func parseHeaders(b *buffer, reader io.Reader) (h.Headers, error) {
	headers := h.Headers{}

	line, err := b.readLine(reader)
	if err != nil {
		return h.Headers{}, err
	}

	for line != "" {
		fmt.Println()
		fmt.Println("Starting iteration...")
		fmt.Printf("Line: %s\n", line)
		fmt.Println("Parsing header")
		headers.Parse(line)
		fmt.Printf("Headers: %v\n", headers)

		fmt.Println("Reading next line")
		line, err = b.readLine(reader)
		if err != nil {
			return h.Headers{}, err
		}
		fmt.Printf("Next line: %s\n", line)
		fmt.Println()
	}

	if len(headers) == 0 {
		return h.Headers{}, fmt.Errorf("empty or malformed headers")
	}

	return headers, nil
}

func containsOnlyCapitalLetters(s string) bool {
	for _, c := range s {
		if !unicode.IsUpper(c) {
			return false
		}
	}

	return true
}
