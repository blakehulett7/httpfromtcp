package main

import (
	"io"
	"strings"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	lines_channel := make(chan string)

	go func() {
		var line string
		buffer := make([]byte, 8)
		for {
			_, err := file.Read(buffer)
			if err == io.EOF {
				file.Close()
				break
			}

			parts := strings.Split(string(buffer), "\n")

			for i, part := range parts {
				if i == len(parts)-1 {
					line += part
					break
				}

				line += part
				lines_channel <- line
				line = ""
			}
		}
		lines_channel <- line
		close(lines_channel)
	}()

	return lines_channel
}
