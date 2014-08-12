//process lines
package main

import "strings"

//takes a line and a user and processes it.
func ProcessLine(user *User, msg string) {
	user.SendLine(msg)
	args := strings.Split(msg, " ")
	checkme := strings.ToLower(args[0])
	switch checkme {
	case "quit":
		user.SendLine("Lol QUIT!!!")
		user.Quit()
	}
}
