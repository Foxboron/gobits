package main

import (
	"os"
	"os/signal"
)

func main() {

	config := get_config("./config")

	// WTF Go
	var chans []Channel
	var nets []Network

	for _, i := range config.Networks {
		for _, n := range i.Channels {
			chans = append(chans, Channel{name: n})
		}
		nets = append(nets, Network{server: i.Server, port: i.Port, channels: chans, nick: config.Nick})
	}
	networks := Networks{servers: nets}

	go networks.ConnectAll()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		println(sig)
		networks.CloseAll()
	}
}
