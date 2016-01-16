package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

type jsonconfig struct {
	Nick     string
	Server   string
	Port     string
	Channels []string
	// For later use
	//Networks NetworksType
}

type NetworksType struct {
	Server   string
	Port     string
	Channels []string
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
		if splitted[1] == "HYPE" {
			write(conn, fmt.Sprintf("PRIVMSG %s :HYYYYPPPEEEE", channel))
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

func main() {

	// Configs

	file, e := ioutil.ReadFile("./config")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var config jsonconfig
	json.Unmarshal(file, &config)
	fmt.Printf("Results: %v\n", config)

	log.SetFlags(log.Lshortfile)
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", config.Server, config.Port), conf)
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
		}
	}()

	connbuf := bufio.NewReader(conn)
	user := fmt.Sprintf("USER %s %s %s :Go FTW", config.Nick, config.Nick, config.Nick)
	write(conn, user)
	nick := fmt.Sprintf("NICK %s", config.Nick)
	write(conn, nick)
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
			joinchannels(conn, config.Channels)
		}
		if s["event"] == "PRIVMSG" {
			docmd(conn, s["msg"], s["channel"])
		}
	}
}
