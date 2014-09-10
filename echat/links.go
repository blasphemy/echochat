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
	args := strings.Split(l, " ")
	if len(args) < 2 {
		link.connection.Close()
		return
	}
	link.name = args[0]
	link.id = args[1]
	link.SendLine(fmt.Sprintf("%s %s", config.ServerName, config.ServerID))
	l, _ = b.ReadString('\n')
	args = strings.Split(l, " ")
	if args[0] != "PW" /* && args[1] != Sha1String(config.LinkPassword) */ {
		log.Printf("Attempted server connection has incorrect password, disconnecting")
		link.SendLine("DIE")
		link.connection.Close()
		return
	}
	links[link.id] = link
	link.SendLine("OK")
	link.SendUsers()
	for {
		l, _ = b.ReadString('\n')
		if config.Debug {
			log.Printf("LINK %s: %s", link.name, l)
		}
		link.route(l)
	}
}

func (link *ServerLink) HandleRequests() {
	//Relevant information has been exchanged (not really, but this could change)
	b := bufio.NewReader(link.connection)
	for {
		//All the magic happens here
		l, _ := b.ReadString('\n')
		if config.Debug {
			log.Printf("LINK %s: %s", link.name, l)
		}
		link.route(l)
	}
}

func (link *ServerLink) route(msg string) {
	args := strings.Split(msg, " ")
	checkme := strings.TrimSpace(strings.ToLower(args[0]))
	switch checkme {
	case "send_to_user":
		link.SendToUserHandler(args)
		break
	case "users":
		link.HandleUsersLine()
		break
	case "sendusers":
		link.SendUsers()
		break
	case "wat":
		break
	default:
		link.SendLine("WAT")
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
		link.SendUsers()
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
	link.SendLine("USERS")
	link.connection.Write(append(lol2, '\n'))
}

func (link *ServerLink) HandleUsersLine() {
	newmap := map[string]*RemoteUser{}
	b := bufio.NewReader(link.connection)
	line, err := b.ReadBytes('\n')
	if err != nil {
		log.Printf("ERROR %s, %s", link.name, err.Error())
		return
	}
	err = json.Unmarshal(line, newmap)
	if err != nil {
		log.Printf("Error unmarshaling users from %s", link.name)
	} else {
		link.users = newmap
	}
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
	new := &User{
		nick:     user.Nick,
		id:       user.Id,
		ip:       user.Ip,
		host:     user.Host,
		realip:   user.Realip,
		realhost: user.Realhost,
		realname: user.Realname,
		remote:   true,
	}
	return new
}