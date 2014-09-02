package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type configuration struct {
	ServerName         string
	ServerDescription  string
	DefaultKickReason  string
	DefaultQuitReason  string
	DefaultPartReason  string
	DefaultKillReason  string
	PingTime           time.Duration
	PingCheckTime      time.Duration
	ResolveHosts       bool
	DefaultCmode       string
	StatTime           time.Duration
	Debug              bool
	Cloaking           bool
	OpersKickable      bool
	Salt               string
	ListenIPs          []string
	ListenPorts        []int
	LogChannels        []string
	Opers              map[string]string
	Privacy            bool
	SystemUserName     string
	SystemJoinChannels bool
	AutoJoin           []string
	Logfile            string
}

func SetupConfig() {
	confile, err := ioutil.ReadFile(conf_file_name)
	if err != nil {
		log.Printf("Error reading config file: " + err.Error())
		SetupConfigDefault()
		os.Exit(1)
	} else {
		err := json.Unmarshal(confile, &config)
		if err != nil {
			log.Printf("Error parsing config file: " + err.Error())
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
		if config.DefaultPartReason == "" {
			config.DefaultPartReason = DefaultConf.DefaultPartReason
		}
		if config.DefaultQuitReason == "" {
			config.DefaultQuitReason = DefaultConf.DefaultQuitReason
		}
		if config.PingTime < 5 || config.PingTime > 500 {
			log.Printf("You have a ridiculous ping time, setting it to the default of %s", DefaultConf.PingTime*time.Second)
			config.PingTime = DefaultConf.PingTime
		}
		if config.PingCheckTime > config.PingTime || config.PingCheckTime < 2 {
			newtime := config.PingTime / 2
			log.Printf("Your ping check time does not make senese, setting it to " + string(newtime*time.Second))
			config.PingCheckTime = newtime
		}
		if config.StatTime < 1 {
			config.StatTime = DefaultConf.StatTime
		}
		SystemUser.nick = config.SystemUserName
		SystemUser.host = config.ServerName
		SystemUser.realhost = config.ServerName
		if config.SystemJoinChannels {
			for _, k := range chanlist {
				SystemUser.JoinHandler([]string{"JOIN", k.name})
			}
		}
	}
}

func SetupConfigDefault() {
	log.Printf("Creating default config file")
	k, err := json.MarshalIndent(DefaultConf, "", "\t")
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	err = ioutil.WriteFile(conf_file_name, k, 0644)
	if err != nil {
		log.Printf("Error writing config file: " + err.Error())
		os.Exit(1)
	}
	log.Printf("Config file created at: " + conf_file_name)
	log.Printf("It is highly recommended you edit this before proceeding...")
}
