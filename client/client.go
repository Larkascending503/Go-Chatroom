package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "set server ip")
	flag.IntVar(&serverPort, "port", 8888, "set server port")
}

func NewClient(ServerIp string, ServerPort int) *Client {
	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		flag:       99,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn

	return client
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1. Broadcast")
	fmt.Println("2. Private")
	fmt.Println("3. Rename")
	fmt.Println("0. Exit")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("Invalid input")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			fmt.Println("Broadcast")

		case 2:
			fmt.Println("Private")

		case 3:
			fmt.Println("Rename")

		}
	}
}
func main() {
	flag.Parse()

	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println(">>>>> Failed to connect server")
		return
	} else {
		fmt.Println(">>>>> Connected to server")
	}

	select {}

}
