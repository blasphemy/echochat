package main

import (
	"time"
)

type configuration struct {
	server_name         string
	default_kick_reason string
	ping_time           time.Duration
	ping_check_time     time.Duration
	resolve_hosts       bool
	default_cmode       string
	stat_time           time.Duration
	debug               bool
	salt                string
	listen_ips          []string
	listen_ports        []string
}

func SetupConfig() {
	config = configuration{}
	config.server_name = "test.net.local"
	config.default_kick_reason = "Your behavior is not conductive of the desired environment."
	config.ping_time = 45
	config.ping_check_time = 20
	config.resolve_hosts = true
	config.default_cmode = "nt"
	config.stat_time = 30
	config.debug = true
	config.salt = "testing"
	config.listen_ips = []string{"0.0.0.0"}
	config.listen_ports = []string{"6667", "6668"}
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
