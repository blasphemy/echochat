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
	realhost   string
	epoch      time.Time
	lastrcv    time.Time
	nextcheck  time.Time
	chanlist   map[string]*Channel
	oper       bool
}

func (user *User) PingChecker() {
	for {
		if user.dead {
			return
		}
		if time.Now().After(user.nextcheck) {
			if user.waiting {
				since := time.Since(user.lastrcv).Seconds()
				user.Quit(fmt.Sprintf("Ping Timeout: %.0f seconds", since))
				return
			} else {
				user.SendLine(fmt.Sprintf("PING :%s", config.ServerName))
				user.waiting = true
				user.nextcheck.Add(config.PingTime * time.Second)
				log.Printf("Sent user %s ping", user.nick)
			}
		}
		time.Sleep(config.PingCheckTime * time.Second)
	}
}

func (user *User) QuitCommandHandler(args []string) {
	var reason string
	if len(args) > 1 {
		args[1] = strings.TrimPrefix(args[1], ":")
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
		delete(k.usermodes, user)
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
	msg := strcat(fmt.Sprintf(":%s %.3d %s ", config.ServerName, numeric, user.nick), fmt.Sprintf(NUM[numeric], args...))
	user.SendLine(msg)
}

func NewUser() *User {
	counter = counter + 1
	user := &User{id: counter, nick: "*"}
	user.chanlist = make(map[string]*Channel)
	user.epoch = time.Now()
	user.lastrcv = time.Now()
	user.nextcheck = time.Now().Add(config.PingTime * time.Second)
	userlist[user.id] = user
	return user
}

func (user *User) SetConn(conn net.Conn) {
	user.connection = conn
	user.ip = GetIpFromConn(conn)
	log.Println("New connection from", user.ip)
	user.realhost = user.ip
	if !config.Cloaking {
		user.host = user.ip
	} else {
		if DetermineConnectionType(user.ip) == "IP4" {
			user.host = CloakIP4(user.ip)
		} else {
			user.host = CloakIP6(user.ip)
		}
	}
	if config.ResolveHosts {
		go user.UserHostLookup()
	}
	go user.PingChecker()
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
		return
	}
	if config.Debug {
		log.Printf("Send to %s: %s", user.nick, msg)
	}
}

func (user *User) SendLinef(msg string, args ...interface{}) {
	user.SendLine(fmt.Sprintf(msg, args...))
}

func (user *User) HandleRequests() {
	b := bufio.NewReader(user.connection)
	for {
		if user.dead {
			return
		}
		line, err := b.ReadString('\n')
		if err != nil {
			log.Println("Error reading:", err.Error())
			user.dead = true
			user.Quit("Error")
			return
		}
		if line == "" {
			user.dead = true
			user.Quit("Error")
			return
		}
		line = strings.TrimSpace(line)
		if config.Debug {
			log.Println("Receive from", fmt.Sprintf("%s:", user.nick), line)
		}
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
	args[4] = strings.TrimPrefix(args[4], ":")
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
	user.FireNumeric(RPL_YOURHOST, config.ServerName, software, softwarev)
	user.FireNumeric(RPL_CREATED, epoch)
	//TODO fire RPL_MYINFO when we actually have enough stuff to do it
	user.FireNumeric(RPL_ISUPPORT, isupport)
	user.FireLusers()
}

func (user *User) UserHostLookup() {
	user.SendLinef(":%s NOTICE %s :*** Looking up your hostname...", config.ServerName, user.nick)
	adds, err := net.LookupAddr(user.ip)
	if err != nil {
		user.SendLinef("%s NOTICE %s :*** Unable to resolve your hostname", config.ServerName, user.nick)
		return
	}
	addstring := adds[0]
	adds, err = net.LookupHost(addstring)
	if err != nil {
		user.SendLinef("%s NOTICE %s :*** Unable to resolve your hostname", config.ServerName, user.nick)
		return
	}
	for _, k := range adds {
		if user.ip == k {
			addstring = strings.TrimSuffix(addstring, ".")
			user.realhost = addstring
			if config.Cloaking {
				user.host = CloakHost(addstring)
			} else {
				user.host = addstring
			}
			user.SendLinef(":%s NOTICE %s :*** Found your hostname", config.ServerName, user.nick)
			return
		}
	}
	user.SendLinef(":%s NOTICE %s :*** Your forward and reverse DNS do not match, ignoring hostname", config.ServerName, user.nick)
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
	if channel.HasMode("A") && !user.oper {
		//TODO fire numeric for this
		return
	}
	if channel.HasUser(user) {
		return //should this silently fail?
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
		channel.SendLinef(":%s PART %s :%s", user.GetHostMask(), channel.name, "Leaving")
		delete(channel.userlist, user.id)
		delete(user.chanlist, channel.name)
		delete(channel.usermodes, user)
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
			if j.HasMode("n") && !user.IsIn(j) {
				user.FireNumeric(ERR_CANNOTSENDTOCHAN, j.name)
				return
			}
			if j.HasMode("m") && j.GetUserPriv(user) < 10 {
				user.FireNumeric(ERR_CANNOTSENDTOCHAN, j.name)
				return
			}
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
	if k.GetUserPriv(user) < 100 && k.HasMode("t") {
		return //doesn't have privs to change channel
		// TODO fire the correct numeric
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
		} else {
			if channel.GetUserPriv(user) < 100 {
				user.FireNumeric(ERR_CHANOPRIVSNEEDED, channel.name)
				return
			}
			s := args[2]
			mode := 0
			mcounter := 0
			var targs []string
			if len(args) > 3 {
				targs = args[3:]
			}
			for _, k := range s {
				switch k {
				case '+':
					mode = 2
					break
				case '-':
					mode = 1
					break
				case 'o', 'v':
					if len(targs)-1 < mcounter {
						user.FireNumeric(ERR_NEEDMOREPARAMS, "MODE")
						break
					}
					target := GetUserByNick(targs[mcounter])
					if target == nil {
						user.FireNumeric(ERR_NOSUCHNICK, args[mcounter])
						mcounter = +1
						break
					}
					if !channel.HasUser(target) {
						user.FireNumeric(ERR_USERNOTINCHANNEL, target.nick, channel.name)
						mcounter = +1
						break
					}
					if mode == 2 {
						channel.SetUmode(target, user, string(k))
						mcounter = +1
						break
					}
					if mode == 1 {
						channel.UnsetUmode(target, user, string(k))
						mcounter = +1
						break
					}
					break
				case 't', 'n', 'm', 'A':
					if mode == 2 {
						channel.SetMode(string(k), user)
					} else if mode == 1 {
						channel.UnsetMode(string(k), user)
					}
					break
				}
			}
		}
	}
}

func (user *User) IsIn(channel *Channel) bool {
	for _, k := range user.chanlist {
		if k == channel {
			return true
		}
	}
	return false
}

func (user *User) PingCmd(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "PING")
		return
	}
	args[1] = strings.TrimPrefix(args[1], ":")
	user.SendLinef(":%s PONG %s :%s", config.ServerName, config.ServerName, args[1])
}

func (user *User) WhoHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "WHO")
		return
	}
	whotype := ChanUserNone(args[1])
	if whotype == 1 {
		//its a channel
		k := GetChannelByName(args[1])
		args[1] = k.name
		for _, j := range k.userlist {
			h := strcat("H", k.GetUserPrefix(j))
			user.FireNumeric(RPL_WHOREPLY, k.name, j.ident, j.host, config.ServerName, j.nick, h, ":0", j.realname)
		}
	} else if whotype == 2 {
		//user
		k := GetUserByNick(args[1])
		args[1] = k.nick
		user.FireNumeric(RPL_WHOREPLY, "*", k.ident, k.host, config.ServerName, k.nick, "H", ":0", k.realname)
	}
	user.FireNumeric(RPL_ENDOFWHO, args[1])
}

func (user *User) KickHandler(args []string) {
	if len(args) < 3 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "KICK")
	}
	channel := GetChannelByName(args[1])
	if channel == nil {
		user.FireNumeric(ERR_NOSUCHCHANNEL, args[1])
		return
	}
	target := GetUserByNick(args[2])
	if target == nil {
		user.FireNumeric(ERR_NOSUCHNICK, args[2])
		return
	}
	if channel.GetUserPriv(user) < 100 {
		user.FireNumeric(ERR_CHANOPRIVSNEEDED, channel.name)
		return
	}
	if !channel.HasUser(target) {
		user.FireNumeric(ERR_USERNOTINCHANNEL, target.nick, channel.name)
		return
	}
	var reason string
	if len(args) > 3 {
		reason = strings.Join(args[3:], " ")
		reason = strings.TrimPrefix(reason, ":")
	} else {
		reason = config.DefaultKickReason
	}
	channel.SendLinef(":%s KICK %s %s :%s", user.GetHostMask(), channel.name, target.nick, reason)
	delete(channel.userlist, target.id)
	delete(target.chanlist, channel.name)
	delete(channel.usermodes, target)
	log.Printf("%s kicked %s from %s", user.nick, target.nick, channel.name)
	channel.ShouldIDie()
}

func (user *User) ListHandler(args []string) {
	user.FireNumeric(RPL_LISTSTART)
	for _, k := range chanlist {
		user.FireNumeric(RPL_LIST, k.name, len(k.userlist), k.topic)
	}
	user.FireNumeric(RPL_LISTEND)
}

func (user *User) OperHandler(args []string) {
	if len(args) < 3 {
		user.CommandNotFound(args)
		return
	}
	if config.Opers[args[1]] == args[2] {
		user.oper = true
		user.FireNumeric(RPL_YOUREOPER)
	} else {
		user.CommandNotFound(args)
	}
}

func (user *User) NamesHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "NAMES")
		return
	}
	channel := GetChannelByName(args[1])
	if channel == nil {
		user.FireNumeric(ERR_NOSUCHCHANNEL, args[1])
		return
	}
	channel.FireNames(user)
}
