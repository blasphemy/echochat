//process lines
package main

import (
	"strings"
	"time"
)

//takes a line and a user and processes it.
func ProcessLine(user *User, msg string) {
	user.lastrcv = time.Now()
	user.nextcheck = time.Now().Add(ping_time * time.Second)
	user.waiting = false
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
		if user.registered {
			user.JoinHandler(args)
		} else {
			user.FireNumeric(ERR_NOTREGISTERED)
		}
		break
	case "privmsg":
		if user.registered {
			user.PrivmsgHandler(args)
		} else {
			user.FireNumeric(ERR_NOTREGISTERED)
		}
		break
	case "pong":
		break //lol nothing
	case "lusers":
		user.FireLusers()
		if user.registered {
			user.FireLusers()
		} else {
			user.FireNumeric(ERR_NOTREGISTERED)
		}
		break
	case "part":
		if user.registered {
			user.PartHandler(args)
		} else {
			user.FireNumeric(ERR_NOTREGISTERED)
		}
		break
	case "topic":
		if user.registered {
			user.TopicHandler(args)
		} else {
			user.FireNumeric(ERR_NOTREGISTERED)
		}
		break
	case "protoctl":
		break
	case "mode":
		if user.registered {
			user.ModeHandler(args)
		} else {
			user.FireNumeric(ERR_NOTREGISTERED)
		}
		user.ModeHandler(args)
		break
	case "ping":
		user.PingCmd(args)
		break
	case "who":
		user.WhoHandler(args)
		break
	default:
		user.CommandNotFound(args)
		break
	}
}
