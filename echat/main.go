package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

func main() {
	SetupConfig()
	runtime.GOMAXPROCS(runtime.NumCPU())
	SetupNumerics()
	SetupSytemUser()
	var listeners []net.Listener
	// Listen for incoming connections.
	for _, LISTENING_IP := range config.ListenIPs {
		for _, LISTENING_PORT := range config.ListenPorts {
			l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", LISTENING_IP, LISTENING_PORT))
			if err != nil {
				log.Println("Error listening:", err.Error())
				os.Exit(1)
			} else {
				listeners = append(listeners, l)
				log.Printf("Listening on %s:%d", LISTENING_IP, LISTENING_PORT)
			}
		}
	}
	// Close the listener when the application closes.
	for _, l := range listeners {
		defer l.Close()
	}
	for _, l := range listeners {
		go listenerthing(l)
	}
	periodicStatusUpdate()
}

func listenerthing(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
		} else {
			user := NewUser()
			user.SetConn(conn)
			checkMaxUsers()
			go user.HandleRequests()
		}
	}
}

func checkMaxUsers() {
	if len(userlist) > maxUsers {
		maxUsers = len(userlist)
	}
}

//periodicStatusUpdate shows information about the ircd every 5 seconds or so,
//as well as updating the max users, and goroutines numbers. Since these are
//only ran every 5 seconds or so, it may not be 100% accurate, but who cares
func periodicStatusUpdate() {
	for {
		log.Printf("Status: %d current users", len(userlist))
		log.Printf("Status: %d current channels", len(chanlist))
		if config.Debug {
			log.Printf("Status: %d current Goroutines", runtime.NumGoroutine())
		}
		time.Sleep(config.StatTime * time.Second)
	}
}
