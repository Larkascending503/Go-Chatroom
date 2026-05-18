package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {
	for {
		msg, ok := <-this.C
		if !ok {
			return
		}
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.Broadcast(this, "login")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.Broadcast(this, "logout")
}

func (this *User) sendMsg(msg string) {
	this.conn.Write([]byte(msg + "\n"))
}

func (this *User) DoMessage(msg string) {
	if msg == "LIST" {
		this.server.mapLock.RLock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": " + "online...\n"
			this.sendMsg(onlineMsg)
		}
	} else {
		this.server.Broadcast(this, msg)
	}
}
