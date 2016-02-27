package main

import (
	"fmt"
	"strings"
)

type cmds func(msg string) cmd
type cmd func(msg string)

type CommandInterface interface {
	DoCMD()
}

type Command struct {
	network Network
	msg     map[string]string
	cmds    map[string]cmd
}

func (c Command) WriteToChannel(msg string) {
	c.network.Write(fmt.Sprintf("PRIVMSG %s :%s", c.msg["channel"], msg))
}

func (c Command) DoCMD() {
	c.Register()
	splitted := strings.Split(c.msg["msg"], " ")
	if splitted[0] == "go" {
		if _, ok := c.cmds[splitted[1]]; ok {
			c.cmds[splitted[1]](c.msg["msg"])
		}
	}
}

func (c Command) addCmd(name string, fn cmd) {
	c.cmds[name] = fn
}

func (c Command) Hype(msg string) {
	c.WriteToChannel("HYYYYYPPPEEEE")
}

func (c Command) Register() {
	c.addCmd("HYPE", c.Hype)
}
