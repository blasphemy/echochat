package main

import oldlog "log"
import "strings"
import "fmt"

type Elog struct {
	//fucking nothing
}

func (elog *Elog) Printf(msg string, args ...interface{}) {
	oldlog.Printf(msg, args...)
	SendLineToLogChannels(fmt.Sprintf(msg, args...))
}

func (elog *Elog) Print(args ...interface{}) {
	oldlog.Print(args...)
	SendLineToLogChannels(fmt.Sprint(args...))
}

func (elog *Elog) Println(args ...interface{}) {
	oldlog.Println(args...)
	SendLineToLogChannels(fmt.Sprint(args...))
}

func SendLineToLogChannels(msg string) {
	msg2 := strings.Split(msg, " ")
	for _, k := range config.LogChannels {
		sender := []string{"PRIVMSG", k}
		sender = append(sender, msg2...)
		SystemUser.PrivmsgHandler(sender)
	}
}
