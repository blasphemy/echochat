package main

import (
	"bytes"
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
