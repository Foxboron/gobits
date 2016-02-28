package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
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

func (c Command) WriteToNotice(msg string) {
	c.network.Write(fmt.Sprintf("NOTICE %s :%s", c.msg["user"], msg))
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

func randLineFromFile(msg string) string {
	file, e := ioutil.ReadFile(msg)
	if e != nil {
		fmt.Printf("File error!")
	}
	lines := strings.Split(string(file), "\n")

	length := len(lines)
	r := rand.Intn(length)
	defined := lines[r]
	return string(defined)
}

func (c Command) Hackers(msg string) {
	c.WriteToChannel(randLineFromFile("./hackers"))
}

func (c Command) AddQuote(msg string) {
	msg = strings.SplitN(msg, " ", 3)[2]
	f, err := os.OpenFile("./quotes", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Print("Read error")
	}

	defer f.Close()

	if _, err = f.WriteString(msg + "\n"); err != nil {
		fmt.Print("Read error")
	}
	c.WriteToNotice("Wrote quote!")
}

func (c Command) ReadQuote(msg string) {
	msg = strings.SplitN(msg, " ", 3)[2]
	slice, err := strconv.Atoi(msg)
	if err != nil {
		fmt.Println("Not and int")
	}

	file, e := ioutil.ReadFile("./quotes")
	if e != nil {
		fmt.Printf("File error!")
	}
	lines := strings.Split(string(file), "\n")

	// Hurrdurr programming
	length := len(lines) - 1
	if slice+1 > length {
		c.WriteToChannel("Don't have that many quotes")
	} else {
		defined := lines[slice]
		c.WriteToChannel(string(defined))
	}
}

func (c Command) Register() {
	c.addCmd("HYPE", c.Hype)
	c.addCmd("hackers", c.Hackers)
	c.addCmd("add-quote", c.AddQuote)
	c.addCmd("read-quote", c.ReadQuote)
}
