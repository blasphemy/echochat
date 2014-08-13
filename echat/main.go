package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

var (
	counter   int
	userlist  map[int]User
	max_users int
	epoch     time.Time
)

func main() {
	epoch = time.Now()
	SetupNumerics()
	userlist = make(map[int]User)
	// Listen for incoming connections.
	l, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)
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

func (user *User) Sync() {
	userlist[user.id] = *user
}

func PeriodicStatusUpdate() {
	for {
		fmt.Println("Status:", len(userlist), "current users")
		fmt.Println("Status:", runtime.NumGoroutine(), "current Goroutines")
		fmt.Println("Status:", counter, "total connections")
		fmt.Println("Status:", max_users, "max users")
		time.Sleep(5 * time.Second)
	}
}
