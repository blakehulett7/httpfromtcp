package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var line string

	bytes := make([]byte, 8)
	for {
		_, err = file.Read(bytes)
		if err == io.EOF {
			break
		}

		parts := strings.Split(string(bytes), "\n")

		for i, part := range parts {
			if i == len(parts)-1 {
				line += part
				break
			}

			fmt.Printf("read: %s\n", line+part)
			line = ""
		}
	}

	fmt.Printf("read: %s\n", line)
}
