package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
)

func main() {
	ColorPrint(Blue, "Dominus Iesus Christus")
	ColorPrint(Blue, "----------------------")
	fmt.Println()

	port := ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	ColorPrint(Magenta, fmt.Sprintf("Listening on port %s", port))
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

			lines_channel := getLinesChannel(connection)
			for line := range lines_channel {
				fmt.Println(line)
			}
		}
	}()

	<-done_chan

	fmt.Println()
	ColorPrint(Blue, "------------------")
	ColorPrint(Blue, "Et Spiritus Sancti")
}
