package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

//Channel...
//represents an irc channel
type Channel struct {
	name      string
	epoch     time.Time
	userlist  map[int]*User
	usermodes map[*User]string
	banlist   map[int]*Ban
	cmodes    string
	topic     string
	topichost string
	topictime int64
}

func (channel *Channel) SetTopic(newtopic string, hostmask string) {
	channel.topic = newtopic
	channel.topichost = hostmask
	channel.topictime = time.Now().Unix()
	channel.SendLinef(":%s TOPIC %s :%s", hostmask, channel.name, newtopic)
}

func NewChannel(newname string) *Channel {
	chann := &Channel{name: newname, epoch: time.Now()}
	chann.userlist = make(map[int]*User)
	chann.usermodes = make(map[*User]string)
	chann.banlist = make(map[int]*Ban)
	chanlist[strings.ToLower(chann.name)] = chann
	chann.cmodes = config.DefaultCmode
	log.Printf("Channel %s created", chann.name)
	return chann
}

func (channel *Channel) len() int {
	k := len(channel.userlist)
	var check bool
	for _, k := range config.LogChannels {
		if channel == GetChannelByName(k) {
			check = true
			break
		}
	}
	if channel.HasUser(SystemUser) && !check {
		k--
	}
	return k
}

func (channel *Channel) JoinUser(user *User) {
	channel.userlist[user.id] = user
	if channel.len() == 1 {
		channel.usermodes[user] = "o"
		if config.SystemJoinChannels {
			SystemUser.JoinHandler([]string{"JOIN", channel.name})
		}
	}
	channel.SendLinef(":%s JOIN %s", user.GetHostMask(), channel.name)
	if len(channel.topic) > 0 {
		channel.FireTopic(user)
	}
	channel.FireNames(user)
}

func (channel *Channel) GetUserPrefix(user *User) string {
	if strings.Contains(channel.usermodes[user], "o") {
		return "@"
	}
	if strings.Contains(channel.usermodes[user], "v") {
		return "+"
	}
	return ""
}

func (channel *Channel) FireTopic(user *User) {
	if len(channel.topic) > 0 {
		user.FireNumeric(RPL_TOPIC, channel.name, channel.topic)
		user.FireNumeric(RPL_TOPICWHOTIME, channel.name, channel.topichost, channel.topictime)
	} else {
		user.FireNumeric(RPL_NOTOPIC, channel.name)
	}
}

func (channel *Channel) FireNames(user *User) {
	var buffer bytes.Buffer
	for _, k := range channel.userlist {
		if buffer.Len()+len(channel.GetUserPrefix(k))+len(user.nick) > 500 {
			user.FireNumeric(RPL_NAMEPLY, channel.name, strings.TrimSpace(buffer.String()))
			buffer.Reset()
		}
		buffer.WriteString(channel.GetUserPrefix(k))
		buffer.WriteString(k.nick)
		buffer.WriteString(" ")
	}
	if buffer.Len() > 1 {
		resp := strings.TrimSpace(buffer.String())
		user.FireNumeric(RPL_NAMEPLY, channel.name, resp)
	}
	user.FireNumeric(RPL_ENDOFNAMES, channel.name)
}

func (channel *Channel) GetUserList() []*User {
	list := []*User{}
	for _, k := range channel.userlist {
		list = append(list, k)
	}
	return list
}

func (channel *Channel) GetUserPriv(user *User) int {
	score := 0
	if strings.Contains(channel.usermodes[user], "o") {
		score += 100
	}
	if strings.Contains(channel.usermodes[user], "v") {
		score += 10
	}
	if user.oper {
		score += 1000
	}
	return score
}

func (channel *Channel) ShouldIDie() {
	if channel.len() < 1 {
		if channel.HasUser(SystemUser) {
			SystemUser.PartHandler([]string{"PART", channel.name})
		}
		delete(chanlist, strings.ToLower(channel.name))
		log.Printf("Channel %s has no users, destroying\n", channel.name)
	}
}

func (channel *Channel) FireModes(user *User) {
	user.FireNumeric(RPL_CHANNELMODEIS, channel.name, channel.cmodes)
	user.FireNumeric(RPL_CREATIONTIME, channel.name, channel.epoch.Unix())
}

func (channel *Channel) HasMode(mode string) bool {
	if strings.Contains(channel.cmodes, mode) {
		return true
	} else {
		return false
	}
}

func (channel *Channel) SetUmode(user *User, changing *User, mode string) {
	if channel.CheckYourPrivlege(changing) {
		return
	}
	if !strings.Contains(channel.usermodes[user], mode) {
		channel.usermodes[user] = channel.usermodes[user] + mode
		channel.SendLinef(":%s MODE %s +%s %s", changing.GetHostMask(), channel.name, mode, user.nick)
	}
}

func (channel *Channel) UnsetUmode(user *User, changing *User, mode string) {
	if channel.CheckYourPrivlege(changing) {
		return
	}
	if strings.Contains(channel.usermodes[user], mode) {
		channel.usermodes[user] = strings.Replace(channel.usermodes[user], mode, "", 1)
		channel.SendLinef(":%s MODE %s -%s %s", changing.GetHostMask(), channel.name, mode, user.nick)
	}
}

func (channel *Channel) SetMode(mode string, changing *User) {
	if channel.CheckYourPrivlege(changing) {
		return
	}
	if !strings.Contains(channel.cmodes, mode) {
		channel.cmodes = channel.cmodes + mode
		channel.SendLinef(":%s MODE %s +%s", changing.GetHostMask(), channel.name, mode)
	}
}

func (channel *Channel) UnsetMode(mode string, changing *User) {
	if channel.CheckYourPrivlege(changing) {
		return
	}
	if strings.Contains(channel.cmodes, mode) {
		channel.cmodes = strings.Replace(channel.cmodes, mode, "", 1)
		channel.SendLinef(":%s MODE %s -%s", changing.GetHostMask(), channel.name, mode)
	}
}

func (channel *Channel) HasUser(user *User) bool {
	if channel.userlist[user.id] == user {
		return true
	} else {
		return false
	}
}

func (channel *Channel) SendLinef(msg string, args ...interface{}) {
	for _, k := range channel.userlist {
		k.SendLine(fmt.Sprintf(msg, args...))
	}
}

func (channel *Channel) CheckYourPrivlege(user *User) bool {
	if channel.GetUserPriv(user) < 100 {
		//SHITLORD!
		user.FireNumeric(ERR_CHANOPRIVSNEEDED, channel.name)
		return true //privlege successfully checked.
	} else {
		return false
	}
}

func (channel *Channel) SetBan(m string, user *User) {
	if CheckIfBanExists(channel, m) {
		return
	}
	hm := user.GetHostMask()
	b := NewBan(m, hm)
	channel.banlist[b.id] = b
	channel.SendLinef(":%s MODE %s +b %s", hm, channel.name, m)
}

func (channel *Channel) UnsetBan(m string, user *User) {

}
