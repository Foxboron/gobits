package main

import ()

func main() {

	register_cmds()

	// Configs
	config := get_config("./config")

	for _, i := range config.Networks {
		connect(config.Nick, i.Server, i.Port, i.Channels)
	}
}
