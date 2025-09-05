package main

import (
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines_channel := make(chan string)

	go func() {
		var line string
		buffer := make([]byte, 8)
		for {
			defer f.Close()
			_, err := f.Read(buffer)
			if err != nil {
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
		fmt.Println("Connection closed")
	}()

	return lines_channel
}
