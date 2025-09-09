package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/blakehulett7/httpfromtcp/internal/request"
)

func main() {
	ColorPrint(Cyan, "Dominus Iesus Christus")
	ColorPrint(Cyan, "----------------------")
	fmt.Println()

	port := ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	ColorPrint(Blue, fmt.Sprintf("Listening on port %s", port))
	fmt.Println()

	done_chan := make(chan os.Signal, 1)
	signal.Notify(done_chan, os.Interrupt)

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			fmt.Println("Connection accepted...")

			r, err := request.RequestFromReader(connection)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
			fmt.Println("Headers:")
			for key, value := range r.Headers {
				fmt.Printf("- %s: %s\n", key, value)
			}
		}
	}()

	<-done_chan

	fmt.Println()
	ColorPrint(Cyan, "------------------")
	ColorPrint(Cyan, "Et Spiritus Sancti")
}
