//process lines
package main

import "strings"

//takes a line and a user and processes it.
func ProcessLine(user *User, msg string) {
	args := strings.Split(msg, " ")
	checkme := strings.ToLower(args[0])
	switch checkme {
	case "quit":
		user.Quit()
		break
	case "nick":
		user.NickHandler(args)
		break
	case "user":
		user.UserHandler(args)
		break
	default:
		user.CommandNotFound(args)
		break
	}
}
