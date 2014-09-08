package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type ServerLink struct {
	connection net.Conn
	users      map[string]*RemoteUser
	name       string
	id         string
}

func SetupLinkListeners() {
	var listeners []net.Listener
	log.Printf("test")
	for _, LISTEN_STRING := range config.LinkAdds {
		log.Printf("Attempting to listen for links at %s", LISTEN_STRING)
		l, err := net.Listen("tcp", LISTEN_STRING)
		if err != nil {
			log.Printf("Error Listening on Linking API: " + err.Error())
		} else {
			listeners = append(listeners, l)
			log.Printf("Listening at %s", LISTEN_STRING)
		}
	}
	for _, l := range listeners {
		StartHandlingLinkConnections(l)
	}
}

func StartHandlingLinkConnections(l net.Listener) {
	defer l.Close()
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
	link.SendLine(fmt.Sprintf("%s %s", config.ServerName, config.ServerID))
	l, _ = b.ReadString('\n')
	if strings.Split(l, " ")[0] != "PW" && strings.Split(l, " ")[1] != Sha1String(config.LinkPassword) {
		log.Printf("Attempted server connection has incorrect password, disconnectiong")
		link.SendLine("DIE")
		link.connection.Close()
		return
	}
	link.SendLine("OK")
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
	log.Printf("Forming outbound link to %s", address)
	conn, _ := net.Dial("tcp", address)
	link := &ServerLink{connection: conn}
	b := bufio.NewReader(conn)
	link.SendLine(fmt.Sprintf("%s %s", config.ServerName, config.ServerID))
	l, _ := b.ReadString('\n')
	link.name = strings.Split(l, " ")[0]
	link.id = strings.Split(l, " ")[1]
	link.SendLine(fmt.Sprint("PW %s", Sha1String(config.LinkPassword)))
	l, _ = b.ReadString('\n')
	if l == "OK" {
		log.Printf("Server %s linked", link.name)
		links[link.id] = link
		link.HandleRequests()
	} else {
		log.Printf("Link to %s failed", address)
		conn.Close()
	}
}

func (link *ServerLink) SendLine(msg string) {
	link.connection.Write([]byte(msg + "\n"))
}

func (link *ServerLink) SendUsers() {
	lol := make(map[string]RemoteUser)
	for _, k := range userlist {
		lol[k.id] = RemoteUser{
			Nick:     k.nick,
			Id:       k.id,
			Ip:       k.ip,
			Realip:   k.realip,
			Host:     k.host,
			Realhost: k.realhost,
			Realname: k.realname,
		}
	}
	lol2, _ := json.Marshal(lol)
	link.SendLine(fmt.Sprintf("USERS %s", string(lol2)))
}

func (link *ServerLink) HandleUsersLine(line string) {
  line = strings.TrimPrefix(line, "USERS ")
  json.Unmarshal([]byte(line), link.users)
}

type RemoteUser struct {
	Nick     string
	Id       string
	Ip       string
	Realip   string
	Host     string
	Realhost string
	Realname string
}

func (user *RemoteUser) ToReal() *User {
  
}
