package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
)

type NetworkInterface interface {
	Connect()
	Close()
	Write()
	Read()
	JoinAll()
	Join()
}

type Channel struct {
	name   string
	joined bool
}

type Network struct {
	connected  bool
	server     string
	channels   []Channel
	port       string
	buffer     bufio.Reader
	connection net.Conn
}

type NetworksInterface interface {
	ConnectAll()
	CloseAll()
}

type Networks struct {
	channels []Network
}

func (n Networks) ConnectAll() {
}

func (n Networks) CloseAll() {
}

func (n Network) Join(channel string) {
}

func (n Network) JoinAll() {
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
	n.connected = true
}

func (n Network) Write(msg_type string, channel string, msg string) {
	n.connection.Write([]byte(fmt.Sprintf("%s %s: %s", msg_type, channel, msg)))
}

func (n Network) Read() []byte {
	connbuf := bufio.NewReader(n.connection)
	str, _, err := connbuf.ReadLine()
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
			go val(conn, msg, channel)
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

func connect(nick string, network string, port string, channels []string, wg sync.WaitGroup) {

	log.SetFlags(log.Lshortfile)
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", network, port), conf)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	// Shit idk
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			println(sig)
			conn.Close()
			defer wg.Done()
		}
	}()

	connbuf := bufio.NewReader(conn)
	user_msg := fmt.Sprintf("USER %s %s %s :Go FTW", nick, nick, nick)
	write(conn, user_msg)
	nick_msg := fmt.Sprintf("NICK %s", nick)
	write(conn, nick_msg)
	for {

		str, _, err := connbuf.ReadLine()

		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		s := parse(string(str))
		fmt.Printf("Got: %v\n", s)
		if s["event"] == "PING" {
			write(conn, fmt.Sprintf("PONG :%s", s["msg"]))
		}
		if s["event"] == "266" {
			joinchannels(conn, channels)
		}
		if s["event"] == "PRIVMSG" {
			docmd(conn, s["msg"], s["channel"])
		}
	}
}
