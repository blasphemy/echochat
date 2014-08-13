package main

import (
	"time"
)

type Channel struct {
	name     string
	epoch    time.Time
	userlist map[int]User
}

func NewChannel(newname string) Channel {
	chann := Channel{name: newname, epoch: time.Now()}
	return chann
}
