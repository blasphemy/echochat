//process lines
package main

import "strings"

//takes a line and a user and processes it.
func ProcessLine(user *User, msg string) {
	user.SendLine(msg)
	args := strings.Split(msg, " ")
	switch args[0] {
	case "QUIT":
		user.SendLine("Lol QUIT!!!")
		user.Quit()
	}
}
