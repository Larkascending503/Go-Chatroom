package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(ServerIp string, ServerPort int) *Client {
	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn

	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>> Failed to connect server")
		return
	} else {
		fmt.Println(">>>>> Connected to server")
	}

	select {}

}
