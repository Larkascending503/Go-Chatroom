package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// Online Map
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// Broadcast chan
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.RLock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	// fmt.Println("Connection established")
	user := NewUser(conn, this)

	user.Online()

	// Monitor user
	isLive := make(chan bool)

	reader := bufio.NewReader(conn)

	// set a timer
	go func() {
		timer := time.NewTimer(time.Second * 10)
		defer timer.Stop()

		for {
			select {
			case <-isLive:
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(time.Second * 10)

			case <-timer.C:
				user.sendMsg("You have been kicked out due to 10 seconds of inactivity.\n")
				conn.Close()
				return
			}
		}
	}()

	// read message
	for {
		msgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			user.Offline()
			return
		}

		msg := string(msgBytes)

		if len(msg) > 0 && msg[len(msg)-1] == '\n' {
			msg = msg[:len(msg)-1]
		}

		user.DoMessage(msg)
		// Any incoming message indicates the user is active
		select {
		case isLive <- true:
		default:
		}
	}
}

// Start server
func (this *Server) Start() {
	// 1. socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// 4. close
	defer listener.Close()

	// 2. accept
	go this.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		// 3. handle
		go this.Handler(conn)
	}
}
