package main

import (
	"errors"
	"fmt"
	"net"
)

type cmd func(conn net.Conn, msg string, channel string)

var _cmds = map[string]cmd{}

func addCmd(name string, fn cmd) {
	_cmds[name] = fn
}

func getCmd(name string) (cmd, error) {
	if val, ok := _cmds[name]; ok {
		return val, nil
	}
	return nil, errors.New("No cmd found")
}

func register_cmds() {
	addCmd("HYPE", func(conn net.Conn, msg string, channel string) {
		write(conn, fmt.Sprintf("PRIVMSG %s :HYYYYPPPEEEE", channel))
	})
}
