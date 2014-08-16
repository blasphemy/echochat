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
		break
	default:
		user.CommandNotFound(args)
		break
	}
}
