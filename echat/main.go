package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "6667"
	CONN_TYPE = "tcp"
)

var (
	counter   int
	userlist  map[int]User
	max_users int
)

func main() {
	SetupNumerics()
	userlist = make(map[int]User)
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)

	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	go PeriodicStatusUpdate()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
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

func AddUserToList(user User) {
	userlist[user.id] = user
}

func PeriodicStatusUpdate() {
	for {
		fmt.Println("Status:", len(userlist), "current users")
    fmt.Println("Status:", runtime.NumGoroutine, "current Goroutines")
		fmt.Println("Status:", counter, "total connections")
		fmt.Println("Status:", max_users, "max users")
		time.Sleep(5 * time.Second)
	}
}
