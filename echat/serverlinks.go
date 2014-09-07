package main

import (
	"net"
)

type ServerLink struct {
	connection net.Conn
}

var (
	links []*ServerLink
)

func SetupLinkListeners() {
	var listeners []net.Listener
	for _, LISTEN_STRING := range config.LinkAddresses {
		l, err := net.Listen("tcp", LISTEN_STRING)
		if err != nil {
			log.Printf("Error Listening on Linking API: " + err.Error())
		} else {
			listeners = append(listeners, l)
		}
	}
	for _, l := range listeners {
		defer l.Close()
	}
	for _, l := range listeners {
		StartHandlingLinkRequests(l)
	}
}

func StartHandlingLinkRequests(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting link connection: " + err.Error())
		} else {
			link := &ServerLink{connection: conn}
			links = append(links, link)
		}
	}
}
