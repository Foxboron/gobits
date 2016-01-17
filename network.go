package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
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

func parse(msg string) map[string]string {
	splitted := strings.SplitN(msg, " :", 3)
	userinfo := strings.Split(splitted[0], " ")
	event := ""
	channel := ""

	fmt.Printf("Internal: %v\n", userinfo)
	if len(userinfo) > 1 {
		event = userinfo[1]
		if len(userinfo) >= 3 {
			channel = userinfo[2]
		}
	} else {
		event = splitted[0]
	}

	info := map[string]string{
		"msg":     splitted[len(splitted)-1],
		"event":   event,
		"channel": channel,
	}
	return info
}

func joinchannels(conn net.Conn, channels []string) {
	for _, i := range channels {
		write(conn, fmt.Sprintf("JOIN :%s", i))
	}
}

func docmd(conn net.Conn, msg string, channel string) {
	splitted := strings.Split(msg, " ")
	if splitted[0] == "go" {
		val, err := getCmd(splitted[1])
		if err != nil {
			println("Found no functions")
		} else {
			val(conn, msg, channel)
		}

	}
}

func write(conn net.Conn, msg string) {
	println("Wrote: ", msg)
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		println("Write to server failed:", err.Error())
	}
}

func connect(nick string, network string, port string, channels []string) {
}
