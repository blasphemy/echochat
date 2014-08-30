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
	counter           = 1
	userlist          = make(map[int]*User)
	chanlist          = make(map[string]*Channel)
	maxUsers          int
	maxRoutines       int
	epoch             = time.Now()
	opercount         = 0
	SystemUser        = &User{
		user:       "system",
		ident:      "system",
		id:         0,
		realname:   "system user",
		userset:    true,
		registered: true,
		ip:         "127.0.0.1",
		realip:     "127.0.0.1",
		epoch:      time.Now(),
		chanlist:   make(map[string]*Channel),
		oper:       true,
		system:     true,
	}
	log = &Elog{}
)
