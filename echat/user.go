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
	dead       bool
	nickset    bool
	connection net.Conn
	id         int
	realname   string
	userset    bool
	registered bool
	ip         string
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
	newip := conn.RemoteAddr().String()
	newip = strings.Split(newip, ":")[0]
	user := User{id: counter, connection: conn, ip: newip}
	user.host = user.ip
	AddUserToList(user)
	return user
}

func (user *User) SendLine(msg string) {
	msg = fmt.Sprintf("%s\n", msg)
	if user.dead {
		return
	}
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
		if line == "" {
			user.Quit()
			break
		}
		fmt.Println("Received Line: ", line)
		ProcessLine(user, line)
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
