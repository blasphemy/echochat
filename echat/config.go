package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	ListenPorts       []int
}

func SetupConfig() {
	confile, err := ioutil.ReadFile("echochat.json")
	if err != nil {
		log.Print("Error reading config file echochat.json: " + err.Error())
		os.Exit(1)
	} else {
		err := json.Unmarshal(confile, &config)
		if err != nil {
			log.Print("Error parsing config file: " + err.Error())
			os.Exit(1)
		}
	}
}

func SetupConfigDefault() {
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
		ListenPorts:       []int{6667, 6668},
	}
	k, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		log.Print(err.Error())
	}
	log.Print(string(k))
	ioutil.WriteFile("echochat.json", k, 0644)
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
