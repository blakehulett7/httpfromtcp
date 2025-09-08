package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
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

			lines_channel := getLinesChannel(connection)
			for line := range lines_channel {
				fmt.Println(line)
			}
		}
	}()

	<-done_chan

	fmt.Println()
	ColorPrint(Cyan, "------------------")
	ColorPrint(Cyan, "Et Spiritus Sancti")
}
