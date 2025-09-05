package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}

	lines_channel := getLinesChannel(file)
	for line := range lines_channel {
		fmt.Printf("read: %s\n", line)
	}
}
