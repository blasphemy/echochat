package main

import "time"

const (
	software       = "echochat"
	softwarev      = "v0.1"
	isupport       = "NAMESX CHANTYPES=#& PREFIX=(ov)@+"
	conf_file_name = "echochat.json"
)

var (
	valid_chan_prefix = []string{"#", "&"}
	global_bad_chars  = []string{":", "!", "@", "*", "(", ")", "<", ">", ",", "~", "/", "\\"}
	config            configuration
	counter           int
	userlist          = make(map[int]*User)
	chanlist          = make(map[string]*Channel)
	maxUsers          int
	maxRoutines       int
	epoch             = time.Now()
)
