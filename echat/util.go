package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
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
	if nick == strings.ToLower(SystemUser.nick) {
		return SystemUser
	}
	for _, k := range userlist {
		if strings.ToLower(k.nick) == nick {
			return k
		}
	}
	return nil
}

func SetUserIPInfo(user *User) {
	ip := user.connection.RemoteAddr().String()
	if !strings.HasPrefix(ip, "[") {
		//ipv4
		user.ConnType = "IP4"
		user.realip = strings.Split(ip, ":")[0]
	} else {
		user.ConnType = "IP6"
		ip = strings.Split(ip, "]")[0]
		ip = strings.TrimPrefix(ip, "[")
		user.realip = ip
	}
}

func GetChannelByName(name string) *Channel {
	return chanlist[strings.ToLower(name)]
}

func SendToMany(msg string, list []*User) {
	users := make(map[*User]int)
	for _, j := range list {
		users[j] = 0
	}
	for j, _ := range users {
		j.SendLine(msg)
	}
}

func ValidChanName(name string) bool {
	if ChanHasBadChars(name) {
		return false
	}
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
	msg = strings.TrimPrefix(msg, ":")
	return msg
}

func NickHasBadChars(nick string) bool {
	for _, k := range global_bad_chars {
		if strings.Contains(nick, k) {
			return true
		}
	}
	for _, k := range valid_chan_prefix {
		if strings.Contains(nick, k) {
			return true
		}
	}
	return false
}

func ChanHasBadChars(nick string) bool {
	for _, k := range global_bad_chars {
		if strings.Contains(nick, k) {
			return true
		}
	}
	return false
}

func ChanUserNone(name string) int {
	if GetChannelByName(name) != nil {
		return 1
	} else if GetUserByNick(name) != nil {
		return 2
	} else {
		return 0
	}
}

func WildcardMatch(text string, pattern string) bool {
	cards := strings.Split(pattern, "*")
	for _, card := range cards {
		index := strings.Index(text, card)
		if index == -1 {
			return false
		}
		text = text[index+len(card):]
	}
	return true
}

func CloakIP6(text string) string {
	stuff := strings.Split(text, ":")
	stuff[len(stuff)-1] = CloakString(stuff[len(stuff)-1], config.Salt)
	stuff[len(stuff)-2] = CloakString(stuff[len(stuff)-2], config.Salt)
	text = strings.Join(stuff, ":")
	return text
}

func CloakIP4(text string) string {
	stuff := strings.Split(text, ".")
	stuff[len(stuff)-1] = CloakString(stuff[len(stuff)-1], config.Salt)
	stuff[len(stuff)-2] = CloakString(stuff[len(stuff)-2], config.Salt)
	text = strings.Join(stuff, ".")
	return text
}

func CloakHost(text string) string {
	//hostname mode
	stuff := strings.Split(text, ".")
	stuff[0] = CloakString(stuff[0], config.Salt)
	text = strings.Join(stuff, ".")
	return text
}

func CloakString(text string, salt string) string {
	var r string
	r = Sha1String(strcat(text, salt))
	for len(r) < len(text) {
		r = strcat(r, Sha1String(r))
	}
	side := true
	for len(r) > len(text) {
		if side {
			r = r[1:]
			side = false
		} else {
			r = r[:len(r)-1]
			side = true
		}
	}
	return r
}

func Sha1String(text string) string {
	data := []byte(text)
	result := fmt.Sprintf("%x", sha1.Sum(data))
	return result
}

func SetupSytemUser() {
	for _, k := range config.LogChannels {
		SystemUser.JoinHandler([]string{"JOIN", k})
		SystemUser.ModeHandler([]string{"MODE", k, "+A"})
	}
}
