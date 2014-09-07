package main

import (
	"bufio"
	"net"
	"strings"
)

type ServerLink struct {
	connection net.Conn
	users      map[string]string
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
		StartHandlingLinkConnections(l)
	}
}

func StartHandlingLinkConnections(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting link connection: " + err.Error())
		} else {
			link := &ServerLink{connection: conn}
			links = append(links, link)
			go link.HandleRequests()
		}
	}
}

func (link *ServerLink) HandleRequests() {
	b := bufio.NewReader(link.connection)
	pw, _ := b.ReadString('n')
	if strings.Split(pw, " ")[0] != "PW" {
		log.Printf("Attempted server connection has incorrect password, disconnectiong")
		link.connection.Close()
		return
	}
}
