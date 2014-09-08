package main

import (
	"bufio"
	"bytes"
	"fmt"
	debuglog "log"
	"net"
	"os"
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
	id         string
	realname   string
	userset    bool
	registered bool
	ip         string
	realip     string
	host       string
	realhost   string
	epoch      time.Time
	lastrcv    time.Time
	nextcheck  time.Time
	chanlist   map[string]*Channel
	oper       bool
	system     bool
	ConnType   string
	remote     bool
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
		reason = strings.Join(args[1:], " ")
	} else {
		reason = config.DefaultQuitReason
	}
	user.Quit(reason)
	if user.oper {
		opercount--
	}
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
	msg := fmt.Sprintf(":%s %.3d %s ", config.ServerName, numeric, user.nick) + fmt.Sprintf(NUM[numeric], args...)
	user.SendLine(msg)
}

func NewUser() *User {
	counter = counter + 1
	user := &User{id: fmt.Sprintf("%s%d", config.ServerID, counter), nick: "*"}
	user.chanlist = make(map[string]*Channel)
	user.epoch = time.Now()
	user.lastrcv = time.Now()
	user.nextcheck = time.Now().Add(config.PingTime * time.Second)
	userlist[user.id] = user
	return user
}

func (user *User) SetConn(conn net.Conn) {
	user.connection = conn
	SetUserIPInfo(user)
	log.Printf("New connection from " + user.realip)
	user.realhost = user.realip
	if !config.Cloaking {
		user.host = user.realip
	} else {
		if user.ConnType == "IP4" {
			k := CloakIP4(user.realip)
			user.host = k
			user.ip = k
		} else {
			k := CloakIP6(user.realip)
			user.host = k
			user.ip = k
		}
	}
	if config.MaxUsers > 0 && len(userlist) > config.MaxUsers {
		user.Quit(fmt.Sprintf("Max amount of connections reached (%d)", config.MaxUsers))
		return
	}
	if config.ResolveHosts {
		go user.UserHostLookup()
	}
	go user.PingChecker()
}

func (user *User) SendLine(msg string) {
	msg = fmt.Sprintf("%s\n", msg)
	if user.dead || user.system {
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
		debuglog.Printf("Send to %s: %s", user.nick, msg)
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
			log.Printf("Error reading: " + err.Error())
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
			debuglog.Println("Receive from", fmt.Sprintf("%s:", user.nick), line)
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
	user.realname = strings.TrimSpace(strings.Join(args[4:], " "))
	user.userset = true
	if !user.registered && user.nickset {
		user.UserRegistrationFinished()
	}
}

func (user *User) UserRegistrationFinished() {
	user.registered = true
	log.Printf("User %d finished registration", user.id)
	user.FireNumeric(RPL_WELCOME, user.nick, user.ident, user.host)
	user.FireNumeric(RPL_YOURHOST, config.ServerName, software, softwarev)
	user.FireNumeric(RPL_CREATED, epoch)
	//TODO fire RPL_MYINFO when we actually have enough stuff to do it
	user.FireNumeric(RPL_ISUPPORT, isupport)
	user.LusersHandler([]string{})
	for _, k := range config.AutoJoin {
		user.JoinHandler([]string{"JOIN", k})
	}
}

func (user *User) UserHostLookup() {
	user.SendLinef(":%s NOTICE %s :*** Looking up your hostname...", config.ServerName, user.nick)
	adds, err := net.LookupAddr(user.realip)
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
		if user.realip == k {
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
		//TODO definitely fire numeric for this
		log.Printf("User %s tried to join %s while +A was set.", user.nick, channel.name)
		return
	}
	if channel.HasUser(user) {
		log.Printf("User %s tried to join %s while already joined.", user.nick, channel.name)
		return //should this silently fail?
	}
	if channel.IsUserBanned(user) && !user.oper {
		user.FireNumeric(RPL_BANNEDFROMCHAN, channel.name)
		log.Printf("User %s tried to join %s while banned.", user.nick, channel.name)
		return
	}
	channel.JoinUser(user)
	user.chanlist[channel.name] = channel
	log.Printf("User %s joined %s", user.nick, channel.name)
}

func (user *User) LusersHandler(args []string) {
	user.FireNumeric(RPL_LUSERCLIENT, len(userlist), 0, 1) //0 services and 1 servers for now
	user.FireNumeric(RPL_LUSEROP, opercount)
	user.FireNumeric(RPL_LUSERUNKNOWN, 0) //also 0...
	user.FireNumeric(RPL_LUSERCHANNELS, len(chanlist))
	user.FireNumeric(RPL_LUSERME, len(userlist), 1)
}

func (user *User) PartHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "PART")
		return
	}
	var reason string
	if len(args) > 2 {
		args[2] = strings.TrimPrefix(args[2], ":")
		reason = strings.Join(args[2:], " ")
	} else {
		reason = config.DefaultPartReason
	}
	channel := GetChannelByName(args[1])
	if channel != nil {
		channel.SendLinef(":%s PART %s :%s", user.GetHostMask(), channel.name, reason)
		delete(channel.userlist, user.id)
		delete(user.chanlist, channel.name)
		delete(channel.usermodes, user)
		log.Printf("User %s PART %s: %s", user.nick, channel.name, reason)
		channel.ShouldIDie()
	} //else?
}

func (user *User) PrivmsgHandler(args []string) {
	if len(args) < 3 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "PRIVMSG")
		return
	}
	//is ValidChanName even needed here anymore?
	if ValidChanName(args[1]) { //TODO part of this should be sent to the channel "object"
		//presumably a channel
		j := GetChannelByName(args[1])
		if j != nil {
			if j.HasMode("n") && !user.IsIn(j) && !user.oper {
				user.FireNumeric(ERR_CANNOTSENDTOCHAN, j.name)
				return
			}
			userpriv := j.GetUserPriv(user)
			if j.HasMode("m") && userpriv < 10 && !user.oper {
				user.FireNumeric(ERR_CANNOTSENDTOCHAN, j.name)
				return
			}
			if j.IsUserBanned(user) && userpriv < 10 && !user.oper {
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
			var logchan bool
			for _, testc := range config.LogChannels {
				if GetChannelByName(testc) == j {
					logchan = true
				}
			}
			if !logchan && !config.Privacy {
				log.Printf("User %s CHANMSG %s: %s", user.nick, j.name, msg)
			}
			return
		} else {
			user.FireNumeric(ERR_NOSUCHCHANNEL, args[1])
			return
		}
	} else {
		//maybe its a user
		target := GetUserByNick(args[1])
		if target != nil {
			msg := FormatMessageArgs(args)
			target.SendLinef(":%s PRIVMSG %s :%s", user.GetHostMask(), target.nick, msg)
			if !config.Privacy {
				log.Printf("User %s PRIVMSG %s: %s", user.nick, target.nick, msg)

			}
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
				case 'b':
					if len(targs)-1 < mcounter {
						channel.FireBanlist(user)
						break
					}
					if mode == 2 {
						channel.SetBan(targs[mcounter], user)
						mcounter++
						break
					}
					if mode == 1 {
						channel.UnsetBan(targs[mcounter], user)
						mcounter++
						break
					}
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
			h := "H" + k.GetUserPrefix(j)
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
	if user.system {
		return //This could be bad.
	}
	if target.oper && !config.OpersKickable {
		return // >:|
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
		opercount++
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

func (user *User) RehashHandler(args []string) {
	if user.oper {
		SetupConfig()
		user.FireNumeric(RPL_REHASHING, conf_file_name)
		log.Printf("OPER %s requested rehash...", user.nick)
	} else {
		user.CommandNotFound(args)
	}
}

func (user *User) ShutdownHandler(args []string) {
	if user.oper {
		log.Printf("Shutdown requested by OPER %s", user.nick)
		for _, k := range userlist {
			k.Quit(fmt.Sprintf("Server is being shutdown by %s", user.nick))
		}
		os.Exit(0)
	} else {
		user.CommandNotFound(args)
	}
}

func (user *User) KillHandler(args []string) {
	if user.oper {
		if len(args) < 2 {
			user.FireNumeric(ERR_NEEDMOREPARAMS, "KILL")
			return
		} else {
			var reason string
			bill := GetUserByNick(args[1])
			if bill == nil {
				user.FireNumeric(ERR_NOSUCHNICK, args[1])
				return
			} else {
				if len(args) > 2 {
					reason = strings.Join(args[2:], " ")
					reason = strings.TrimPrefix(reason, ":")
				} else {
					reason = config.DefaultKillReason
				}
				bill.Quit(fmt.Sprintf("KILL: %s", reason))
				log.Printf("oper %s killed %s (%s)", user.nick, bill.nick, reason)
			}
		}
	} else {
		user.CommandNotFound(args)
	}
}

func (user *User) WhoisHandler(args []string) {
	if len(args) < 2 {
		user.FireNumeric(ERR_NEEDMOREPARAMS, "WHOIS")
		return
	}
	target := GetUserByNick(args[1])
	if target == nil {
		user.FireNumeric(ERR_NOSUCHNICK, args[1])
		return
	}
	var buf bytes.Buffer
	for _, k := range target.chanlist {
		buf.WriteString(k.name + " ")
	}
	chanstring := strings.TrimSpace(buf.String())
	user.FireNumeric(RPL_WHOISUSER, target.nick, target.ident, target.host, target.realname)
	user.FireNumeric(RPL_WHOISCHANNELS, target.nick, chanstring)
	user.FireNumeric(RPL_WHOISSERVER, target.nick, config.ServerName, config.ServerDescription)
	if target.oper {
		user.FireNumeric(RPL_WHOISOPERATOR, target.nick)
	}
	if user.oper || user == target {
		user.FireNumeric(RPL_WHOISHOST, target.nick, target.realhost, target.realip)
	} else {
		user.FireNumeric(RPL_WHOISHOST, target.nick, target.host, target.ip)
	}
	user.FireNumeric(RPL_ENDOFWHOIS, target.nick)
}
