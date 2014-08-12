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
	dead       bool
	connection net.Conn
}

func (user *User) Quit() {
	user.dead = true
	user.connection.Close()
}

//set a user's nick. Does not do any error checking
func (user *User) SetNick(new string) {
	user.nick = new
}

//returns a user's nick in string format
func (user *User) GetNick() string {
	return user.nick
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
		if user.dead {
			break
		}
		line, err := b.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			user.Quit()
		}
		line = strings.TrimSpace(line)
		fmt.Println("Received Line: ", line)
		// Send a response back to person contacting us.
		go ProcessLine(user, line)
	}
}
