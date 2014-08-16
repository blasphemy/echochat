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

func GetUserByNick(nick string) *User {
	nick = strings.ToLower(nick)
	for _, k := range userlist {
		if strings.ToLower(k.nick) == nick {
			return k
		}
	}
	return nil
}

func GetIpFromConn(conn net.Conn) string {
	ip := conn.RemoteAddr().String()
	ip = strings.Split(ip, ":")[0]
	return ip
}

func GetChannelByName(name string) (bool, *Channel) {
	for _, k := range chanlist {
		if strings.ToLower(name) == strings.ToLower(k.name) {
			return true, k
		}
	}
	channel := NewChannel(name)
	return false, channel
}

func SendToMany(msg string, list []*User) {
	list2 := []*User{}
	for _, j := range list {
		match := false
		for _, k := range list2 {
			if j == k {
				match = true
				break
			}
		}
		if !match {
			list2 = append(list2, j)
		}
	}
	for _, j := range list2 {
		j.SendLine(msg)
	}
}

func ValidChanName(name string) bool {
	for _, k := range valid_chan_prefix {
		if strings.HasPrefix(name, k) {
			return true
		}
	}
	return false
}

//IMPORTANT: args must ABSOLUTELY be a valid privmsg command, or this will not work
//validity does not depend on leading ":", I don't care that much
func FormatMessageArgs(args []string) string {
	msg := strings.Join(args[2:], " ")
	msg = StripLeading(msg, ":")
	return msg
}

//removes the first instance of toremove, only if it has it
//r can be more than one char, but I wouldn't recommend it
func StripLeading(msg string, r string) string {
	if strings.HasPrefix(msg, r) {
		msg = strings.Replace(msg, r, "", 1)
	}
	return msg
}

func NickHasBadChars(nick string) bool {
	for _, k := range forbidden_nick_chars {
		if strings.Contains(nick, k) {
			return true
		}
	}
	return false
}
