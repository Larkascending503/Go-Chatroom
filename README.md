# Go Chatroom

A simple TCP-based chatroom system built with Go, supporting multiple concurrent users with broadcast messaging, private messaging, and more.

## Features

- **Broadcast Messaging** — Send messages to all online users
- **Private Messaging** — Send messages to a specific user
- **Online User List** — View all currently connected users
- **Rename** — Change your display name
- **Auto-Kick** — Users are automatically disconnected after 60 seconds of inactivity

## Project Structure

```
communication-system/
├── server/
│   ├── main.go       # Server entry point
│   ├── server.go     # Server core logic (TCP listener, broadcast, connection handler)
│   └── user.go       # User struct and message handling (LIST, RENAME, PVT)
├── client/
│   └── client.go     # Client with interactive menu
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.26+

### Start the Server

```bash
cd server
go run .
```

The server starts on `127.0.0.1:8888` by default.

### Start a Client

```bash
cd client
go run . -ip 127.0.0.1 -port 8888
```

#### Command-line Flags

| Flag    | Default     | Description       |
| ------- | ----------- | ----------------- |
| `-ip`   | `127.0.0.1` | Server IP address |
| `-port` | `8888`      | Server port       |

### Client Menu

Once connected, the client presents an interactive menu:

```
1. Broadcast  — Send a message to all users
2. Private    — Send a private message to a specific user
3. Rename     — Change your username
0. Exit       — Disconnect and quit
```

## Usage Examples

**Broadcast:** Select option `1`, type your message, and press Enter. Type `exit` to return to the menu.

**Private Message:** Select option `2`, view the online user list, enter the target username and your message. Type `exit` to return to the menu.

**Rename:** Select option `3` and enter your new username.

## License

MIT