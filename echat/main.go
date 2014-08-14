package main

import (
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

var (
	counter   int
	userlist  map[int]*User
	chanlist  map[string]*Channel
	max_users int
	epoch     time.Time
)

func main() {
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
	go PeriodicStatusUpdate()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)

		}
		// Handle connections in a new goroutine.
		user := NewUser(conn)
		go CheckMaxUsers()
		go user.HandleRequests()
	}

}

func CheckMaxUsers() {
	if len(userlist) > max_users {
		max_users = len(userlist)
	}
}

func PeriodicStatusUpdate() {
	for {
		log.Println("Status:", len(userlist), "current users")
		log.Println("Status:", len(chanlist), "current channels")
		log.Println("Status:", runtime.NumGoroutine(), "current Goroutines")
		time.Sleep(5 * time.Second)

	}
}
