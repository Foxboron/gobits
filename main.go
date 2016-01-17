package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup
	register_cmds()

	// Configs
	config := get_config("./config")

	for _, i := range config.Networks {
		wg.Add(1)
		go connect(config.Nick, i.Server, i.Port, i.Channels, wg)
	}
	wg.Wait()
}
