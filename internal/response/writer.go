package response

import (
	"fmt"
	"io"

	h "github.com/blakehulett7/httpfromtcp/internal/headers"
)

type HandlerError struct {
	StatusCode StatusCode
	Error      error
}

type writerState int

const (
	Waiting writerState = iota
	StatusLineWritten
	HeadersWritten
	BodyWritten
	Done
)

type Writer struct {
	writer  io.Writer
	Headers h.Headers
	state   writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  Waiting,
	}
}

func (w *Writer) WriteStatusLine(status_code StatusCode) error {
	if w.state != Waiting {
		return fmt.Errorf("WriteStatusLine must be called first...")
	}

	err := WriteStatusLine(w.writer, status_code)
	if err != nil {
		return err
	}

	w.state = StatusLineWritten
	return nil
}

func (w *Writer) WriteHeaders() error {
	if w.state != StatusLineWritten && w.state != BodyWritten {
		return fmt.Errorf("WriteHeaders must be called directly after WriteStatusLine")
	}

	err := WriteHeaders(w.writer, w.Headers)
	if err != nil {
		return err
	}

	w.state = HeadersWritten
	return nil
}

func (w *Writer) WriteBody(data []byte) (int, error) {
	if w.state != HeadersWritten {
		return 0, fmt.Errorf("WriteBody must be called directly after WriteHeaders")
	}

	bytes_written, err := w.writer.Write(data)
	if err != nil {
		return bytes_written, err
	}

	w.state = BodyWritten

	return bytes_written, nil
}

func (w *Writer) WriteError(err HandlerError, content_type string) {
	w.WriteStatusLine(StatusCode(err.StatusCode))
	message := []byte(err.Error.Error())
	w.Headers = GetDefaultHeaders(len(message))
	w.Headers.Set("Content-Type", content_type)
	w.WriteHeaders()
	w.writer.Write(message)
}

func (w *Writer) WriteChunkedBody(data []byte) (int, error) {
	w.writer.Write(fmt.Appendf(nil, "%x\r\n", len(data)))
	w.writer.Write(data)
	return w.writer.Write([]byte("\r\n"))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	w.state = BodyWritten
	return w.writer.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(h h.Headers) error {
	w.Headers = h
	w.WriteHeaders()
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
