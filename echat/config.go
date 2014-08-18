package main

const (
	sname           = "test.net.local"
	software        = "echochat"
	softwarev       = "v0.1"
	CONN_HOST       = "0.0.0.0"
	CONN_PORT       = "6667"
	ping_time       = 45 //something
	ping_check_time = 20 // time between the user's ping checks
	resolvehosts    = false //Note: make forward confirmed reverse dns optional.
	isupport        = "NAMESX CHANTYPES=#& PREFIX=(o)@"
)

var (
	valid_chan_prefix = []string{"#", "&"}
	global_bad_chars  = []string{":", "!", "@", "*", "(", ")", "<", ">", ",", "~", "/", "\\"}
)
