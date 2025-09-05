package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes := make([]byte, 8)
	for {
		_, err = file.Read(bytes)
		if err == io.EOF {
			break
		}

		fmt.Printf("read: %s\n", bytes)
	}
}
