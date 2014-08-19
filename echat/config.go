package main

const (
	sname               = "test.net.local"
	software            = "echochat"
	softwarev           = "v0.1"
	CONN_HOST           = "0.0.0.0"
	CONN_PORT           = "6667"
	default_kick_reason = "Your behavior is not conductive of the desired environment."
	ping_time           = 45   //something
	ping_check_time     = 20   // time between the user's ping checks
	resolvehosts        = true //Note: make forward confirmed reverse dns optional.
	isupport            = "NAMESX CHANTYPES=#& PREFIX=(ov)@+"
	default_cmode       = "nt"
	stattime            = 30
	debug               = false
)

var (
	valid_chan_prefix = []string{"#", "&"}
	global_bad_chars  = []string{":", "!", "@", "*", "(", ")", "<", ">", ",", "~", "/", "\\"}
)
