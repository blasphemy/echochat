//process lines
package main

//takes a line and a user and processes it.
func ProcessLine(user *User, msg string) {
	user.SendLine(msg)
}
