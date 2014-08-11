package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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

func (user *User) SendLine(msg string) {
	conn := user.GetConn()
	msg = fmt.Sprintf("%s\n", msg)
	conn.Write([]byte(msg))
}

func (user *User) HandleRequests() {
	conn := user.GetConn()
	b := bufio.NewReader(conn)
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())

		}
		line = strings.TrimSpace(line)
		fmt.Println("Received Line: ", line)
		// Send a response back to person contacting us.
		go ProcessLine(user, line)
		// Close the connection when you're done with it.

	}
	conn.Close()
}
