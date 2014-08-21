package main

import (
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

var (
	counter     int
	userlist    map[int]*User
	chanlist    map[string]*Channel
	maxUsers    int
	maxRoutines int
	epoch       time.Time
)

func main() {
	epoch = time.Now()
	SetupNumerics()
	userlist = make(map[int]*User)
	chanlist = make(map[string]*Channel)
	var listeners []net.Listener
	// Listen for incoming connections.
	for _, LISTENING_IP := range listen_ips {
		for _, LISTENING_PORT := range listen_ports {
			l, err := net.Listen("tcp", LISTENING_IP+":"+LISTENING_PORT)
			if err != nil {
				log.Println("Error listening:", err.Error())
				os.Exit(1)
			} else {
				listeners = append(listeners, l)
				log.Println("Listening on " + LISTENING_IP + ":" + LISTENING_PORT)
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
		checkMaxUsers()
		gor := runtime.NumGoroutine()
		if gor > maxRoutines {
			maxRoutines = gor
		}
		log.Println("Status:", len(userlist), "current users")
		log.Println("Status:", len(chanlist), "current channels")
		log.Println("Status:", gor, "current Goroutines")
		log.Println("Status:", maxRoutines, "max Goroutines")
		time.Sleep(stattime * time.Second)

	}
}
