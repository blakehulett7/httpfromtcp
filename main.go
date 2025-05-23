package main

import (
	"fmt"
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

	channel := getLinesChannel(messages)
	for line := range channel {
		fmt.Printf("read: %s\n", line)
	}
}
