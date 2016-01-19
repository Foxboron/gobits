package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type NetworkInterface interface {
	Connect()
	Close()
	Write(msg string)
	Read() byte
	JoinAll()
	Join(channel string)
}

type Channel struct {
	name   string
	joined bool
}

type Network struct {
	nick       string
	connected  bool
	server     string
	channels   []Channel
	port       string
	buffer     *bufio.Reader
	connection net.Conn
}

type NetworksInterface interface {
	ConnectAll()
	CloseAll()
}

type Networks struct {
	servers []Network
}

func (n Networks) ConnectAll() {
	for _, i := range n.servers {
		go i.Connect()
	}
}

func (n Networks) CloseAll() {
	for _, i := range n.servers {
		i.Close()
	}
}

func (n Network) Join(channel Channel) {
	n.Write(fmt.Sprintf("JOIN :%s", channel.name))
	channel.joined = true
}

func (n Network) JoinAll() {
	for _, i := range n.channels {
		n.Join(i)
	}
}

func (n Network) Connect() {
	log.SetFlags(log.Lshortfile)
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", n.server, n.port), conf)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	n.connection = conn
	n.buffer = bufio.NewReader(n.connection)
	n.connected = true

	user_msg := fmt.Sprintf("USER %s %s %s :Go FTW", n.nick, n.nick, n.nick)
	n.Write(user_msg)
	nick_msg := fmt.Sprintf("NICK %s", n.nick)
	n.Write(nick_msg)

	for {
		str := n.Read()
		s := parse(string(str))
		fmt.Printf("Got: %v\n", s)
		if s["event"] == "PING" {
			n.Write(fmt.Sprintf("PONG :%s", s["msg"]))
		}
		if s["event"] == "266" {
			n.JoinAll()
		}
		if s["event"] == "PRIVMSG" {
			go docmd(n, s)
		}
	}
}

func (n Network) Write(msg string) {
	n.connection.Write([]byte(msg + "\n"))
}
func (n Network) Read() []byte {
	str, _, err := n.buffer.ReadLine()
	if err != nil {
		println("Read from server failed:", err.Error())
		n.Close()
	}
	return str
}

func (n Network) Close() {
	n.connection.Close()
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

func docmd(n Network, m map[string]string) {
	splitted := strings.Split(m["msg"], " ")
	if splitted[0] == "go" {
		val, err := getCmd(splitted[1])
		if err != nil {
			println("Found no functions")
		} else {
			msg := val(m["msg"])
			n.Write(fmt.Sprintf("PRIVMSG %s :%s", m["channel"], msg))
		}

	}
}
