package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type ServerLink struct {
	connection net.Conn
	users      map[string]string
	name       string
	id         string
}

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
			go link.Registration()
		}
	}
}

func (link *ServerLink) Registration() {
	b := bufio.NewReader(link.connection)
	l, _ := b.ReadString('\n')
	link.name = strings.Split(l, " ")[0]
	link.id = strings.Split(l, " ")[1]
	links[link.id] = link
	link.connection.Write([]byte(fmt.Sprintf("%s %s", config.ServerName, config.ServerID)))
	l, _ = b.ReadString('\n')
	if strings.Split(l, " ")[0] != "PW" {
		log.Printf("Attempted server connection has incorrect password, disconnectiong")
		link.connection.Close()
		return
	}
	link.connection.Write([]byte("OK"))
	link.HandleRequests()
}

func (link *ServerLink) HandleRequests() {
	//Relevant information has been exchanged (not really, but this could change)
	b := bufio.NewReader(link.connection)

	for {
		//All the magic happens here
		l, _ := b.ReadString('\n')
		link.route(l)
	}
}

func (link *ServerLink) route(msg string) {
	args := strings.Split(msg, " ")
	switch strings.ToLower(args[0]) {
	case "SEND_TO_USER":
		link.SendToUserHandler(args)
		break
	default:
		break
	}
}

//SEND_TO_USER USERID STUFF STUFF STUFF
func (link *ServerLink) SendToUserHandler(args []string) {
	user := userlist[args[1]]
	user.SendLine(strings.Join(args[2:], " "))
}

//TODO error checking
func FormOutgoingLink(address string) {
	conn, _ := net.Dial("tcp", address)
	link := &ServerLink{connection: conn}
	b := bufio.NewReader(conn)
	conn.Write([]byte(fmt.Sprintf("%s %s", config.ServerName, config.ServerID)))
	l, _ := b.ReadString('\n')
	link.name = strings.Split(l, " ")[0]
	link.id = strings.Split(l, " ")[1]
	conn.Write([]byte(fmt.Sprintf("PW %s", config.LinkPassword)))
	l, _ = b.ReadString('\n')
	if l == "OK" {
		links[link.id] = link
		link.HandleRequests()
	} else {
		conn.Close()
	}
}
