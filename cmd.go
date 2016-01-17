package main

import (
	"fmt"
	"net"
)

func register_cmds() {
	addCmd("HYPE", func(conn net.Conn, msg string, channel string) {
		write(conn, fmt.Sprintf("PRIVMSG %s :HYYYYPPPEEEE", channel))
	})
}
