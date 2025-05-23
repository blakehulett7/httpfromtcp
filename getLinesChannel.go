package main

import (
	"bytes"
	"fmt"
	"io"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	channel := make(chan string)
	var line string

	go func() {
		for {
			buffer := make([]byte, 8)
			_, err := file.Read(buffer)

			if err == io.EOF {
				file.Close()
				if line != "" {
					channel <- line
				}
				close(channel)
				return
			}

			if err != nil {
				fmt.Printf("can't read file, err: %v\n", err)
			}

			slice := bytes.Split(buffer, []byte("\n"))

			for i := 0; i < len(slice)-1; i++ {
				part := slice[i]
				channel <- line + string(part)
				line = ""
			}

			part := slice[len(slice)-1]
			line += string(part)
		}
	}()

	return channel
}
