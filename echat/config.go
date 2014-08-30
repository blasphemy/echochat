package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type configuration struct {
	ServerName        string
	ServerDescription string
	DefaultKickReason string
	DefaultQuitReason string
	DefaultPartReason string
	DefaultKillReason string
	PingTime          time.Duration
	PingCheckTime     time.Duration
	ResolveHosts      bool
	DefaultCmode      string
	StatTime          time.Duration
	Debug             bool
	Cloaking          bool
	OpersKickable     bool
	Salt              string
	ListenIPs         []string
	ListenPorts       []int
	LogChannels       []string
	Opers             map[string]string
	Privacy           bool
	SystemUserName    string
}

func SetupConfig() {
	confile, err := ioutil.ReadFile(conf_file_name)
	if err != nil {
		log.Print("Error reading config file: " + err.Error())
		SetupConfigDefault()
		os.Exit(1)
	} else {
		err := json.Unmarshal(confile, &config)
		if err != nil {
			log.Print("Error parsing config file: " + err.Error())
			os.Exit(1)
		}
		if config.SystemUserName == "" {
			config.SystemUserName = DefaultConf.SystemUserName
		}
		if config.ServerName == "" {
			config.ServerName = DefaultConf.ServerName
		}
		if config.ServerDescription == "" {
			config.ServerDescription = DefaultConf.ServerDescription
		}
		if config.DefaultKickReason == "" {
			config.DefaultKickReason = DefaultConf.DefaultKickReason
		}
		if config.DefaultKillReason == "" {
			config.DefaultKillReason = DefaultConf.DefaultKillReason
		}
		SystemUser.nick = config.SystemUserName
		SystemUser.host = config.ServerName
		SystemUser.realhost = config.ServerName
	}
}

func SetupConfigDefault() {
	log.Print("Creating default config file")
	k, err := json.MarshalIndent(DefaultConf, "", "\t")
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
	err = ioutil.WriteFile(conf_file_name, k, 0644)
	if err != nil {
		log.Print("Error writing config file: " + err.Error())
		os.Exit(1)
	}
	log.Print("Config file created at: " + conf_file_name)
	log.Print("It is highly recommended you edit this before proceeding...")
}
