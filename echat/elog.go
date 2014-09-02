package main

import oldlog "log"
import "strings"
import "fmt"

type Elog struct {
	//nothing
}

func (elog *Elog) Printf(msg string, args ...interface{}) {
	oldlog.Printf(msg, args...)
	WriteToLogFile(msg, args...)
	SendLineToLogChannels(fmt.Sprintf(msg, args...))
}

func SendLineToLogChannels(msg string) {
	msg2 := strings.Split(msg, " ")
	for _, k := range config.LogChannels {
		sender := []string{"PRIVMSG", k}
		sender = append(sender, msg2...)
		SystemUser.PrivmsgHandler(sender)
	}
}

func WriteToLogFile(msg string, args ...interface{}) {
	if config.Logfile != "" && LoggingFile != nil {
		_, err := LoggingFile.WriteString(fmt.Sprintf(msg+"\n", args...))
		if err != nil {
			config.Logfile = ""
			log.Printf("ERROR: %s", err.Error())
			log.Printf("Error writing to Logfile %s, disabling file logging", config.Logfile)
		}
	}
}
