package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(connection net.Conn) <-chan string {
	lines_channel := make(chan string)

	go func() {
		var line string
		buffer := make([]byte, 8)
		for {
			_, err := connection.Read(buffer)
			if err == io.EOF {
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
		connection.Close()
		fmt.Println("Connection closed")
	}()

	return lines_channel
}
