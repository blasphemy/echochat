package main

import (
	"bufio"
	"bytes"
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
	nickset    bool
	connection net.Conn
	id         int
	realname   string
	userset    bool
	registered bool
	host       string
}

func (user *User) Quit() {
	user.dead = true
	if user.connection != nil {
		user.connection.Close()
	}
	delete(userlist, user.id)
}

func (user *User) FireNumeric(numeric int, args ...interface{}) {
	msg := strcat(fmt.Sprintf(":%s %.3d ", "test.net.local", numeric), fmt.Sprintf(NUM[numeric], args...))
	user.SendLine(msg)
}

func NewUser(conn net.Conn) User {
	counter = counter + 1
	user := User{id: counter, connection: conn}
	user.host = "lol"
	AddUserToList(user)
	return user
}

func (user *User) SendLine(msg string) {
	msg = fmt.Sprintf("%s\n", msg)
	user.connection.Write([]byte(msg))
}

func (user *User) HandleRequests() {
	b := bufio.NewReader(user.connection)
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
func (user *User) NickHandler(args []string) {
	if CheckNickCollision(args[1]) != false {
		return //TODO handle properly
	}
	if !user.nickset {
		user.nickset = true
	}
	user.nick = args[1]
	fmt.Println("User changed name to", args[1])
	if !user.registered && user.userset {
		user.UserRegistrationFinished()
	}
}

func (user *User) UserHandler(args []string) {
	if len(args) < 5 {
		//ERR_NEEDMOREPARAMS
		return
	}
	user.ident = args[1]
	if strings.HasPrefix(args[4], ":") {
		args[4] = strings.Replace(args[4], ":", "", 1)
	}
	var buffer bytes.Buffer
	for i := 4; i < len(args); i++ {
		buffer.WriteString(args[i])
		buffer.WriteString(" ")
	}
	user.realname = strings.TrimSpace(buffer.String())
	user.userset = true
	if !user.registered && user.nickset {
		user.UserRegistrationFinished()
	}
}

func (user *User) UserRegistrationFinished() {
	user.registered = true
	fmt.Printf("User %d finished registration\n", user.id)
	user.FireNumeric(RPL_WELCOME, user.nick, user.ident, user.host)
}
