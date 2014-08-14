package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type Channel struct {
	name      string
	epoch     time.Time
	userlist  map[int]*User
	topic     string
	topichost string
	topictime int64
}

func NewChannel(newname string) *Channel {
	chann := &Channel{name: newname, epoch: time.Now()}
	chann.userlist = make(map[int]*User)
	chanlist[chann.name] = chann
	return chann
}

func (channel *Channel) JoinUser(user *User) {
	channel.userlist[user.id] = user
	SendToMany(fmt.Sprintf(":%s JOIN %s", user.GetHostMask(), channel.name), channel.GetUserList())
	if len(channel.topic) < 1 {
		user.FireNumeric(RPL_TOPIC, channel.name, channel.topic)
		user.FireNumeric(RPL_TOPICWHOTIME, channel.name, channel.topichost, channel.topictime)
	}
	channel.FireNames(user)
}

func (channel *Channel) FireNames(user *User) {
	var buffer bytes.Buffer
	for _, k := range userlist {
		buffer.WriteString(k.nick)
		buffer.WriteString(" ")
	}
	resp := strings.TrimSpace(buffer.String())
	user.FireNumeric(RPL_NAMEPLY, channel.name, resp)
	user.FireNumeric(RPL_ENDOFNAMES, channel.name)
}

func (channel *Channel) GetUserList() []*User {
	list := []*User{}
	for _, k := range channel.userlist {
		list = append(list, k)
	}
	return list
}
