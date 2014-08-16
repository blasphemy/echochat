package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"
)

//Channel...
//represents an irc channel
type Channel struct {
	name      string
	epoch     time.Time
	userlist  map[int]*User
	topic     string
	topichost string
	topictime int64
}

func (channel *Channel) SetTopic(newtopic string, hostmask string) {
	channel.topic = newtopic
	channel.topichost = hostmask
	channel.topictime = time.Now().Unix()
	for _, k := range channel.userlist {
		k.SendLine(fmt.Sprintf(":%s TOPIC %s :%s", hostmask, channel.name, newtopic))
	}
}

func NewChannel(newname string) *Channel {
	chann := &Channel{name: newname, epoch: time.Now()}
	chann.userlist = make(map[int]*User)
	chanlist[chann.name] = chann
	log.Printf("Channel %s created\n", chann.name)
	return chann
}

func (channel *Channel) JoinUser(user *User) {
	channel.userlist[user.id] = user
	SendToMany(fmt.Sprintf(":%s JOIN %s", user.GetHostMask(), channel.name), channel.GetUserList())
	if len(channel.topic) > 0 {
		channel.FireTopic(user)
	}
	channel.FireNames(user)
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
	for _, k := range userlist {
		if buffer.Len()+len(k.nick) > 500 {
			user.FireNumeric(RPL_NAMEPLY, channel.name, strings.TrimSpace(buffer.String()))
			buffer.Reset()
		}
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

func (channel *Channel) ShouldIDie() {
	if len(channel.userlist) < 1 {
		delete(chanlist, channel.name)
		log.Printf("Channel %s has no users, destroying\n", channel.name)
	}
}
