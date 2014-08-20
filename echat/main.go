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
	log.Print(CloakHost(sname))
	log.Print(CloakIP4("127.0.0.1"))
	log.Print(CloakIP6("2604:180::d8f:be0e"))
	epoch = time.Now()
	SetupNumerics()
	userlist = make(map[int]*User)
	chanlist = make(map[string]*Channel)
	// Listen for incoming connections.
	l, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)

	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	go periodicStatusUpdate()
	for {
		// Listen for an incoming connection.
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
