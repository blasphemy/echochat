package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "6667"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)

	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)

		}
		// Handle connections in a new goroutine.
		user := User{connection: conn}
		go user.HandleRequests()
	}

}
