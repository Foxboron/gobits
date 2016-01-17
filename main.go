package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {

	addCmd("HYPE", func(conn net.Conn, msg string, channel string) {
		write(conn, fmt.Sprintf("PRIVMSG %s :HYYYYPPPEEEE", channel))
	})

	// Configs
	config := get_config("./config")
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
