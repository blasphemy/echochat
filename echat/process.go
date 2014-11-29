//process lines
package main

import (
	"strings"
	"time"
)

type Handler func([]string)

var (
  gpc int64
)

//takes a line and a user and processes it.
func ProcessLine(user *User, msg string) {
	user.lastrcv = time.Now()
	user.Lock()
	user.nextcheck = time.Now().Add(config.PingTime * time.Second)
	user.waiting = false
	user.Unlock()
	args := strings.Split(msg, " ")
	checkme := strings.ToLower(args[0])
	switch checkme {
	case "quit":
		user.QuitCommandHandler(args)
		break
	case "nick":
		user.NickHandler(args)
		break
	case "user":
		user.UserHandler(args)
		break
	case "join":
		user.FireIfRegistered(user.JoinHandler, args)
		break
	case "privmsg":
		user.FireIfRegistered(user.PrivmsgHandler, args)
		break
	case "pong":
		break //lol nothing
	case "lusers":
		user.FireIfRegistered(user.LusersHandler, args)
		break
	case "part":
		user.FireIfRegistered(user.PartHandler, args)
		break
	case "topic":
		user.FireIfRegistered(user.TopicHandler, args)
		break
	case "protoctl":
		break
	case "mode":
		user.FireIfRegistered(user.ModeHandler, args)
		break
	case "ping":
		user.PingCmd(args)
		break
	case "who":
		user.FireIfRegistered(user.WhoHandler, args)
		break
	case "kick":
		user.FireIfRegistered(user.KickHandler, args)
		break
	case "list":
		user.FireIfRegistered(user.ListHandler, args)
		break
	case "names":
		user.FireIfRegistered(user.NamesHandler, args)
		break
	case "oper":
		user.FireIfRegistered(user.OperHandler, args)
		break
	case "rehash":
		user.FireIfRegistered(user.RehashHandler, args)
		break
	case "shutdown":
		user.FireIfRegistered(user.ShutdownHandler, args)
		break
	case "kill":
		user.FireIfRegistered(user.KillHandler, args)
		break
	case "whois":
		user.FireIfRegistered(user.WhoisHandler, args)
		break
	default:
		user.CommandNotFound(args)
		break
	}
}

func (user *User) FireIfRegistered(handler Handler, args []string) {
	if user.registered {
		handler(args)
	} else {
		user.FireNumeric(ERR_NOTREGISTERED)
	}
}
