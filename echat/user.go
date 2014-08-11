package main

import (
	"bufio"
	"fmt"
	"net"
)

type User struct {
	nick       string
	user       string
	ident      string
	ip         string
	connection net.Conn
}

func (user *User) SetConn(in net.Conn) {
	user.connection = in
}

func (user *User) GetConn() net.Conn {
	return user.connection
}

func (user *User) HandleRequests() {
	conn := user.GetConn()
	b := bufio.NewReader(conn)
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())

		}
		fmt.Println("Received Line: ", line)
		// Send a response back to person contacting us.
		conn.Write([]byte("Message received."))
		// Close the connection when you're done with it.

	}
	conn.Close()
}
