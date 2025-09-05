package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println(err)
		}
	}
}
