package main

import (
	"crypto/tls"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net"
	"os"
	"runtime"
	"time"
)

func main() {
	SetupNumerics()
	SetupConfig()
	SetupPool()
	var listeners []net.Listener
	// Listen for incoming connections.
	var tconfig tls.Config
	var cert tls.Certificate
	var tlser error
	if len(config.TLSPorts) > 0 {
		cert, tlser = tls.LoadX509KeyPair(config.TLSCertPath, config.TLSKeyPath)
		if tlser == nil {
			tconfig = tls.Config{Certificates: []tls.Certificate{cert}}
		} else {
			log.Printf("TLS ERR: %s", tlser.Error())
		}
	}
	for _, LISTENING_IP := range config.ListenIPs {
		for _, LISTENING_PORT := range config.ListenPorts {
			l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", LISTENING_IP, LISTENING_PORT))
			if err != nil {
				log.Printf("Error listening: " + err.Error())
				os.Exit(1)
			} else {
				listeners = append(listeners, l)
				log.Printf("Listening on %s:%d", LISTENING_IP, LISTENING_PORT)
			}
		}
		if len(config.TLSPorts) > 0 && tlser == nil {
			for _, LISTENING_PORT := range config.TLSPorts {
				l, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", LISTENING_IP, LISTENING_PORT), &tconfig)
				if err != nil {
					log.Printf("Error listening: " + err.Error())
					os.Exit(1)
				} else {
					listeners = append(listeners, l)
					log.Printf("TLS Listening on %s:%d", LISTENING_IP, LISTENING_PORT)
				}
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
			log.Printf("Error accepting: " + err.Error())
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

func SetupPool() {
	RedisPool = &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort))
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", config.RedisPassword); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}
