package main

var (
	NUM map[int]string
)

func SetupNumerics() {
	NUM = make(map[int]string)
	NUM[RPL_WELCOME] = "Welcome to the Internet Relay Network %s!%s@%s"
	NUM[RPL_YOURHOST] = "Your host is %s, running %s version %s"
	NUM[RPL_CREATED] = "This server was created %s"
	NUM[RPL_MYINFO] = "%s %s %s %s"
}

const (
	RPL_WELCOME  = 001
	RPL_YOURHOST = 002
	RPL_CREATED  = 003
	RPL_MYINFO   = 004
)
