package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var servAddr string = "irc.velox.pw:6666"
var botnick string = "Gobits"

func parse(msg string) []string {
	splitted := strings.Split(msg, ":")
	return splitted
}

func write(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		println("Write to server failed:", err.Error())
	}
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	connbuf := bufio.NewReader(conn)
	user := fmt.Sprintf("USER %s %s %s :Go FTW", botnick, botnick, botnick)
	write(conn, user)
	nick := fmt.Sprintf("NICK %s", botnick)
	write(conn, nick)
	for {

		str, _, err := connbuf.ReadLine()

		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		s := parse(string(str))
		fmt.Printf("%v", s)
	}

}
