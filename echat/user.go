package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
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
	chanlist   map[string]*Channel
}

func (user *User) QuitCommandHandler(args []string) {
	var reason string
	if len(args) > 1 {
		args[1] = StripLeading(args[1], ":")
		var buffer bytes.Buffer
		for i := 1; i < len(args); i++ {
			buffer.WriteString(args[i])
			buffer.WriteString(" ")
		}
		reason = strings.TrimSpace(buffer.String())
	} else {
		reason = "Leaving"
	}
	user.Quit(reason)
}

func (user *User) Quit(reason string) {
	targets := []*User{user}
	for _, k := range user.chanlist {
		targets = append(targets, k.GetUserList()...)
		delete(k.userlist, user.id)
		delete(user.chanlist, k.name)
		k.ShouldIDie()
	}
	SendToMany(fmt.Sprintf(":%s QUIT :%s", user.GetHostMask(), reason), targets)
	user.SendLine(fmt.Sprintf("ERROR :Closing Link: %s (%s)", user.host, reason))
	user.dead = true
	if user.connection != nil {
		user.connection.Close()
	}
	delete(userlist, user.id)
}

func (user *User) FireNumeric(numeric int, args ...interface{}) {
	msg := strcat(fmt.Sprintf(":%s %.3d %s ", sname, numeric, user.nick), fmt.Sprintf(NUM[numeric], args...))
	user.SendLine(msg)
}

func NewUser(conn net.Conn) *User {
	userip := GetIpFromConn(conn)
	log.Println("New connection from", userip)
	counter = counter + 1
	user := &User{id: counter, connection: conn, ip: userip, nick: "*"}
	user.chanlist = make(map[string]*Channel)
	user.host = user.ip
	user.epoch = time.Now()
	userlist[user.id] = user
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
		user.dead = true
		user.Quit("Error")
		log.Printf("Error sending message to %s, disconnecting\n", user.nick)
	}
	log.Printf("Send to %s: %s", user.nick, msg)
}

func (user *User) HandleRequests() {
	b := bufio.NewReader(user.connection)
	for {
		if user.dead {
			break
		}
		line, err := b.ReadString('\n')
		if err != nil {
			log.Println("Error reading:", err.Error())
			user.dead = true
			user.Quit("Error")
		}
		if line == "" {
			user.dead = true
			user.Quit("Error")
			break
		}
		line = strings.TrimSpace(line)
		log.Println("Receive from", fmt.Sprintf("%s:", user.nick), line)
		go ProcessLine(user, line)
	}
}
func (user *User) NickHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NONICKNAMEGIVEN)
		return
	}
	if GetUserByNick(args[1]) != nil {
		user.FireNumeric(ERR_NICKNAMEINUSE, args[1])
		return
	}
	if !user.nickset {
		user.nickset = true
	} else if user.registered {
		targets := []*User{}
		targets = append(targets, user)
		for _, k := range user.chanlist {
			targets = append(targets, k.GetUserList()...)
		}
		SendToMany(fmt.Sprintf(":%s NICK %s", user.GetHostMask(), args[1]), targets)
	}
	user.nick = args[1]
	if !user.registered && user.userset {
		user.UserRegistrationFinished()
	}
}

func (user *User) UserHandler(args []string) {
	if len(args) < 5 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "USER")
		return
	}
	user.ident = args[1]
	args[4] = StripLeading(args[4], ":")
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
	log.Printf("User %d finished registration\n", user.id)
	user.FireNumeric(RPL_WELCOME, user.nick, user.ident, user.host)
	user.FireNumeric(RPL_YOURHOST, sname, software, softwarev)
	user.FireNumeric(RPL_CREATED, epoch)
	//TODO fire RPL_MYINFO when we actually have enough stuff to do it
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
}

func (user *User) CommandNotFound(args []string) {
	user.FireNumeric(ERR_UNKNOWNCOMMAND, args[0])
}

func (user *User) GetHostMask() string {
	return fmt.Sprintf("%s!%s@%s", user.nick, user.ident, user.host)
}

func (user *User) JoinHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "JOIN")
		return
	}
	if !ValidChanName(args[1]) {
		user.FireNumeric(ERR_NOSUCHCHANNEL, args[1])
		return
	}
	_, channel := GetChannelByName(args[1])
	channel.JoinUser(user)
	user.chanlist[channel.name] = channel
}

func (user *User) PrivmsgHandler(args []string) {
	if len(args) < 3 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "PRIVMSG")
		return
	}
	if ValidChanName(args[1]) { //TODO part of this should be sent to the channel "object"
		//presumably a channel
		k, j := GetChannelByName(args[1])
		if k {
			//channel exists, send the message
			msg := FormatMessageArgs(args)
			list := j.GetUserList()
			for _, l := range list {
				if l != user {
					l.SendLine(fmt.Sprintf(":%s PRIVMSG %s :%s", user.GetHostMask(), j.name, msg))
				}
			}
			return
		} else {
			//channel didnt exist but get channel by name makes one anyways, lets kill it...
			user.FireNumeric(ERR_NOSUCHCHANNEL, args[1])
			j.ShouldIDie()
			return
		}
	} else {
		//maybe its a user
		target := GetUserByNick(args[1])
		if target != nil {
			msg := FormatMessageArgs(args)
			target.SendLine(fmt.Sprint(":%s PRIVMSG %s :%s", user.GetHostMask(), target.nick, msg))
		}
	}
}
