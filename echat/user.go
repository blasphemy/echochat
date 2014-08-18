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
	waiting    bool
	connection net.Conn
	id         int
	realname   string
	userset    bool
	registered bool
	ip         string
	host       string
	epoch      time.Time
	lastrcv    time.Time
	nextcheck  time.Time
	chanlist   map[string]*Channel
}

func (user *User) PingChecker() {
	for {
		if user.dead {
			break
		}
		if time.Now().After(user.nextcheck) {
			if user.waiting {
				since := time.Since(user.lastrcv).Seconds()
				user.Quit(fmt.Sprintf("Ping Timeout: %.0f seconds", since))
				break
			} else {
				user.SendLine(fmt.Sprintf("PING :%s", sname))
				user.waiting = true
				user.nextcheck.Add(ping_time * time.Second)
				log.Printf("Sent user %s ping", user.nick)
			}
		}
		time.Sleep(ping_check_time * time.Second)
	}
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
	log.Printf("User %s Quit (%s)", user.nick, reason)
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
	user.SendLinef("ERROR :Closing Link: %s (%s)", user.host, reason)
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
	user.lastrcv = time.Now()
	user.nextcheck = time.Now().Add(ping_time * time.Second)
	userlist[user.id] = user
	if resolvehosts {
		go user.UserHostLookup()
	}
	go user.PingChecker()
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
	//log.Printf("Send to %s: %s", user.nick, msg)
}

func (user *User) SendLinef(msg string, args ...interface{}) {
	user.SendLine(fmt.Sprintf(msg, args...))
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
		//log.Println("Receive from", fmt.Sprintf("%s:", user.nick), line)
		ProcessLine(user, line)
	}
}
func (user *User) NickHandler(args []string) {
	oldnick := user.nick
	if len(args) < 2 {
		user.FireNumeric(ERR_NONICKNAMEGIVEN)
		return
	}
	if NickHasBadChars(args[1]) {
		user.FireNumeric(ERR_ERRONEOUSNICKNAME, args[1])
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
	log.Printf("User %s changed nick to %s", oldnick, user.nick)
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
	user.FireNumeric(RPL_ISUPPORT, isupport)
	user.FireLusers()
}

func (user *User) UserHostLookup() {
	user.SendLinef(":%s NOTICE %s :*** Looking up your hostname...", sname, user.nick)
	adds, err := net.LookupAddr(user.ip)
	if err != nil {
		user.SendLinef("%s NOTICE %s :*** Unable to resolve your hostname", sname, user.nick)
		return
	}
	addstring := adds[0]
	adds, err = net.LookupHost(addstring)
	if err != nil {
		user.SendLinef("%s NOTICE %s :*** Unable to resolve your hostname", sname, user.nick)
		return
	}
	for _, k := range adds {
		if user.ip == k {
			user.host = addstring
			user.SendLinef(":%s NOTICE %s :*** Found your hostname", sname, user.nick)
			return
		}
	}
	user.SendLinef(":%s NOTICE %s :*** Your forward and reverse DNS do not match, ignoring hostname", sname, user.nick)
}

func (user *User) CommandNotFound(args []string) {
	log.Printf("User %s attempted unknown command %s", user.nick, args[0])
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
	channel := GetChannelByName(args[1])
	if channel == nil {
		channel = NewChannel(args[1])
	}
	channel.JoinUser(user)
	user.chanlist[channel.name] = channel
	log.Printf("User %s joined %s", user.nick, channel.name)
}

func (user *User) FireLusers() {
	user.FireNumeric(RPL_LUSERCLIENT, len(userlist), 0, 1) //0 services and 1 servers for now
	user.FireNumeric(RPL_LUSEROP, 0)                       //also 0 for now
	user.FireNumeric(RPL_LUSERUNKNOWN, 0)                  //also 0...
	user.FireNumeric(RPL_LUSERCHANNELS, len(chanlist))
	user.FireNumeric(RPL_LUSERME, len(userlist), 1)
}

func (user *User) PartHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "PART")
		return
	}
	channel := GetChannelByName(args[1])
	if channel != nil {
		list := channel.GetUserList()
		delete(channel.userlist, user.id)
		delete(user.chanlist, channel.name)
		SendToMany2f(list, ":%s PART %s :%s", user.GetHostMask(), channel.name, "Leaving")
		log.Printf("User %s PART %s: %s", user.nick, channel.name, "Leaving")
		channel.ShouldIDie()
	} //else?
}

func (user *User) PrivmsgHandler(args []string) {
	if len(args) < 3 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "PRIVMSG")
		return
	}
	if ValidChanName(args[1]) { //TODO part of this should be sent to the channel "object"
		//presumably a channel
		j := GetChannelByName(args[1])
		if j != nil {
			//channel exists, send the message
			msg := FormatMessageArgs(args)
			list := j.GetUserList()
			for _, l := range list {
				if l != user {
					l.SendLinef(":%s PRIVMSG %s :%s", user.GetHostMask(), j.name, msg)
				}
			}
			log.Printf("User %s CHANMSG %s: %s", user.nick, j.name, msg)
			return
		} else {
			//channel didnt exist but get channel by name makes one anyways, lets kill it...
			user.FireNumeric(ERR_NOSUCHCHANNEL, args[1])
			return
		}
	} else {
		//maybe its a user
		target := GetUserByNick(args[1])
		if target != nil {
			msg := FormatMessageArgs(args)
			target.SendLinef(":%s PRIVMSG %s :%s", user.GetHostMask(), target.nick, msg)
			log.Printf("User %s PRIVMSG %s: %s", user.nick, target.nick, msg)
		}
	}
}

func (user *User) TopicHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "TOPIC")
		return
	}
	k := GetChannelByName(args[1])
	if k == nil {
		user.FireNumeric(ERR_NOSUCHCHANNEL, args[1])
		return
	}
	if len(args) < 3 {
		k.FireTopic(user)
		return
	}
	msg := FormatMessageArgs(args)
	k.SetTopic(msg, user.GetHostMask())
}

func (user *User) ModeHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "MODE")
		return
	}
	if ChanUserNone(args[1]) == 1 {
		channel := GetChannelByName(args[1])
		if len(args) < 3 {
			//just digging around...
			channel.FireModes(user)
			log.Printf("User %s requested modes for %s", user.nick, channel.name)
		}
	}
}
