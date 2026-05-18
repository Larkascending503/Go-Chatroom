package main

import (
	"fmt"
	"net"
	"sync"
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
	user := NewUser(conn)

	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	this.Broadcast(user, "login")

	select {}
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
