package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>> Please input new name:")
	fmt.Scanln(&client.Name)
	sendMsg := "RENAME|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (client *Client) Broadcast() {
	var charMsg string
	fmt.Println(">>>>> Please input message(exit to quit):")
	fmt.Scanln(&charMsg)

	for charMsg != "exit" {
		if len(charMsg) != 0 {
			sendMsg := charMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		charMsg = ""
		fmt.Println(">>>>> Please input message(exit to quit):")
		fmt.Scanln(&charMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "LIST\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}
func (client *Client) Private() {
	client.SelectUsers()
	fmt.Println(">>>>> Please input user name(exit to quit):")

	var remoteName string
	var msg string
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>> Please input message(exit to quit):")
		fmt.Scanln(&msg)

		for msg != "exit" {
			if msg != "" {
				sendMsg := "PVT|" + remoteName + "|" + msg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err:", err)
					return
				}
			}

			msg = ""
			fmt.Println(">>>>> Please input message(exit to quit):")
			fmt.Scanln(&msg)
		}

		remoteName = ""
		fmt.Println(">>>>> Please input user name(exit to quit):")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) DealResponer() {
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			client.Broadcast()
		case 2:
			client.Private()
		case 3:
			client.UpdateName()
		case 0:
			fmt.Println("Exit")
			return
		}
	}
}
func main() {
	flag.Parse()

	client := NewClient(serverIP, serverPort)

	client.DealResponer()

	if client == nil {
		fmt.Println(">>>>> Failed to connect server")
		return
	} else {
		fmt.Println(">>>>> Connected to server")
	}

	client.Run()

}
