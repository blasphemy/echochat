package main

import (
	"encoding/json"
	"log"
	"time"
)

type configuration struct {
	ServerName        string
	DefaultKickReason string
	PingTime          time.Duration
	PingCheckTime     time.Duration
	ResolveHosts      bool
	DefaultCmode      string
	StatTime          time.Duration
	Debug             bool
	Salt              string
	ListenIps         []string
	ListenPorts       []string
}

func SetupConfig() {
	config = configuration{
		ServerName:        "test.net.local",
		DefaultKickReason: "Your behavior is not conductive of the desired environment.",
		PingTime:          45,
		PingCheckTime:     20,
		ResolveHosts:      true,
		DefaultCmode:      "nt",
		StatTime:          30,
		Debug:             true,
		Salt:              "testing",
		ListenIps:         []string{"0.0.0.0"},
		ListenPorts:       []string{"6667", "6668"},
	}
	k, err := json.Marshal(config)
	if err != nil {
		log.Print(err.Error())
	}
	log.Print(string(k))
}

const (
	software  = "echochat"
	softwarev = "v0.1"
	isupport  = "NAMESX CHANTYPES=#& PREFIX=(ov)@+"
)

var (
	valid_chan_prefix = []string{"#", "&"}
	global_bad_chars  = []string{":", "!", "@", "*", "(", ")", "<", ">", ",", "~", "/", "\\"}
	config            configuration
)
