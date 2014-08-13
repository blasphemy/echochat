package main

import (
	"bytes"
	"net"
	"strings"
)

func strcat(s1 string, s2 string) string {
	var buffer bytes.Buffer
	buffer.WriteString(s1)
	buffer.WriteString(s2)
	return buffer.String()
}

func CheckNickCollision(nick string) bool {
	nick = strings.ToLower(nick)
	for _, k := range userlist {
		if strings.ToLower(k.nick) == nick {
			return true
		}
	}
	return false
}

func GetIpFromConn(conn net.Conn) string {
	ip := conn.RemoteAddr().String()
	ip = strings.Split(ip, ":")[0]
	return ip
}
