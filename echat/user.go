package main

import "net"

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
