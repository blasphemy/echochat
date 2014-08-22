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
	Cloaking          bool
	Salt              string
	ListenIps         []string
	ListenPorts       []int
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
	}
}

func SetupConfigDefault() {
	log.Print("Creating default config file")
	config = configuration{
		ServerName:        "test.net.local",
		DefaultKickReason: "Your behavior is not conductive of the desired environment.",
		PingTime:          45,
		PingCheckTime:     20,
		ResolveHosts:      true,
		DefaultCmode:      "nt",
		StatTime:          30,
		Debug:             false,
		Cloaking:          false,
		Salt:              "default",
		ListenIps:         []string{"0.0.0.0"},
		ListenPorts:       []int{6667, 6668, 6669},
	}
	k, err := json.MarshalIndent(config, "", "\t")
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
