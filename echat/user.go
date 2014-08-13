package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
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
	epoch      time.Time
}

func (user *User) Quit() {
	user.dead = true
	user.Sync()
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
	userip := GetIpFromConn(conn)
	fmt.Println("New connection from", userip)
	counter = counter + 1
	user := User{id: counter, connection: conn, ip: userip, nick: "*"}
	user.host = user.ip
	user.epoch = time.Now()
	user.Sync()
	go user.UserHostLookup()
	return user
}

func (user *User) SendLine(msg string) {
	msg = fmt.Sprintf("%s\n", msg)
	if user.dead {
		return
	}
	_, err := user.connection.Write([]byte(msg))
	if err != nil {
		user.Quit()
		fmt.Printf("Error sending message to %s, disconnecting\n", user.nick)
	}
	fmt.Printf("Send to %s: %s", user.nick, msg)
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
		if line == "" {
			user.Quit()
			break
		}
		line = strings.TrimSpace(line)
		fmt.Println("Receive from", fmt.Sprintf("%s:", user.nick), line)
		ProcessLine(user, line)
	}
}
func (user *User) NickHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NONICKNAMEGIVEN)
		return
	}
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
	user.Sync()
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
	user.Sync()
	if !user.registered && user.nickset {
		user.UserRegistrationFinished()
	}
}

func (user *User) UserRegistrationFinished() {
	user.registered = true
	fmt.Printf("User %d finished registration\n", user.id)
	user.FireNumeric(RPL_WELCOME, user.nick, user.ident, user.host)
	user.FireNumeric(RPL_YOURHOST, sname, software, softwarev)
	user.FireNumeric(RPL_CREATED, epoch)
	//TODO fire RPL_MYINFO when we actually have enough stuff to do it
	user.Sync()
}

func (user *User) UserHostLookup() {
	user.SendLine(fmt.Sprintf(":%s NOTICE %s :*** Looking up your hostname...", sname, user.nick))
	adds, err := net.LookupAddr(user.ip)
	if err != nil {
		user.SendLine(fmt.Sprintf("%s NOTICE %s :*** Unable to resolve your hostname", sname, user.nick))
		return
	}
	addstring := adds[0]
	adds, err = net.LookupHost(addstring)
	if err != nil {
		user.SendLine(fmt.Sprintf("%s NOTICE %s :*** Unable to resolve your hostname", sname, user.nick))
		return
	}
	for _, k := range adds {
		if user.ip == k {
			user.host = addstring
			user.SendLine(fmt.Sprintf(":%s NOTICE %s :*** Found your hostname", sname, user.nick))
			return
		}
	}
	user.SendLine(fmt.Sprintf(":%s NOTICE %s :*** Your forward and reverse DNS do not match, ignoring hostname", sname, user.nick))
	user.Sync()
}

func (user *User) Sync() {
	userlist[user.id] = *user
}

func (user *User) CommandNotFound(args []string) {
	user.FireNumeric(ERR_UNKNOWNCOMMAND, args[0])
}
