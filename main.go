package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	const file_path = "./messages.txt"

	messages, err := os.Open(file_path)
	if err != nil {
		fmt.Printf("can't open %v, err: %v\n", file_path, err)
		return
	}
	defer messages.Close()

	var line string
	for {
		buffer := make([]byte, 8)
		_, err := messages.Read(buffer)

		if err == io.EOF {
			if line != "" {
				fmt.Printf("read: %v\n", line)
			}
			return
		}

		if err != nil {
			fmt.Printf("can't read file, err: %v\n", err)
		}

		slice := bytes.Split(buffer, []byte("\n"))

		for i := 0; i < len(slice)-1; i++ {
			part := slice[i]
			fmt.Printf("read: %v\n", line+string(part))
			line = ""
		}

		part := slice[len(slice)-1]
		line += string(part)
	}
}
