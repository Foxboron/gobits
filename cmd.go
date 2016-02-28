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
	help    map[string]string
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
			msg := strings.SplitN(c.msg["msg"], " ", 3)
			if len(msg) >= 3 {
				c.cmds[splitted[1]](msg[2])
			} else {
				c.cmds[splitted[1]]("")
			}
		}
	}
}

func (c Command) addCmd(name string, help string, fn cmd) {
	c.cmds[name] = fn
	c.help[name] = help
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
	f, err := os.OpenFile("./quotes", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Print("Read error")
	}

	defer f.Close()

	if _, err = f.WriteString("\n" + msg); err != nil {
		fmt.Print("Read error")
	}
	c.WriteToNotice("Wrote quote!")
}

func (c Command) ReadQuote(msg string) {
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

func (c Command) Deadpool(msg string) {
	c.WriteToChannel("https://i.imgur.com/CJfu35j.gif")
}

func (c Command) HeheJPG(msg string) {
	c.WriteToChannel("https://iskrembilen.com/hehe.jpg")
}

func (c Command) HeheGIF(msg string) {
	c.WriteToChannel("https://iskrembilen.com/hehe.gif")
}

func (c Command) HehePNG(msg string) {
	c.WriteToChannel("https://iskrembilen.com/hehe.png")
}

func (c Command) Help(msg string) {
	ret := ""
	if msg == "" {
		for k := range c.help {
			ret += k + " "
		}
	} else {
		ret = msg + ": " + c.help[msg]
	}
	c.WriteToChannel(ret)
}

func (c Command) Register() {
	c.addCmd("HYPE", "GO LANG HYPE!", c.Hype)
	c.addCmd("hackers", "Awesomesauce quotes", c.Hackers)
	c.addCmd("add-quote", "Add a quote", c.AddQuote)
	c.addCmd("read-quote", "Read a quote", c.ReadQuote)
	c.addCmd("no", "Yeah no....", c.Deadpool)
	c.addCmd("hehe-jpg", "ehehehehehe", c.HeheJPG)
	c.addCmd("hehe-gif", "ehehehehe", c.HeheGIF)
	c.addCmd("hehe-png", "eheheheh", c.HehePNG)
	c.addCmd("help", "Get help!", c.Help)
}
