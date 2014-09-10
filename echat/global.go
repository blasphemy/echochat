package main

import (
	"os"
	"time"
)

const (
	software       = "echochat"
	softwarev      = "v0.1"
	isupport       = "NAMESX CHANTYPES=#& PREFIX=(ov)@+"
	conf_file_name = "echochat.json"
)

var (
	StartupIncomplete = true //used to determine if the ircd is up and running yet
	valid_chan_prefix = []string{"#", "&"}
	global_bad_chars  = []string{":", "!", "@", "*", "(", ")", "<", ">", ",", "~", "/", "\\"}
	config            *configuration
	counter           = 1
	userlist          = make(map[string]*User)
	chanlist          = make(map[string]*Channel)
	maxUsers          int
	maxRoutines       int
	epoch             = time.Now()
	opercount         = 0
	SystemUser        = &User{
		user:       "system",
		ident:      "system",
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
	DefaultConf = &configuration{
		ServerName:         "test.net.local",
		ServerDescription:  "A test server",
		DefaultKickReason:  "Your behavior is not conductive of the desired environment.",
		DefaultKillReason:  "Your behavior is not conductive of the desired environment.",
		DefaultQuitReason:  "Leaving",
		DefaultPartReason:  "Leaving",
		PingTime:           45,
		PingCheckTime:      20,
		ResolveHosts:       true,
		DefaultCmode:       "nt",
		StatTime:           30,
		Debug:              false,
		Cloaking:           false,
		OpersKickable:      false,
		Salt:               "default",
		ListenIPs:          []string{"0.0.0.0"},
		ListenPorts:        []int{6667, 6668, 6669},
		LogChannels:        []string{"#log", "#opers"},
		Opers:              map[string]string{"default": "password"},
		Privacy:            false,
		SystemUserName:     "system",
		AutoJoin:           []string{"#default"},
		SystemJoinChannels: false,
		Logfile:            "echochat.log",
		ServerID:           "A",
		MaxUsers:           0,
		LinkPassword:       "secure",
	}
	log         = &Elog{}
	LoggingFile *os.File
	links       = make(map[string]*ServerLink)
)
