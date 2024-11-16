package main

import (
	"fmt"
	"net"
	"os"
)

func readUDP(connection *net.UDPConn, bufferSize int) (data []byte, remoteAddr *net.UDPAddr) {
	data = make([]byte, bufferSize)
	_, remoteAddr, err := connection.ReadFromUDP(data)

	//log.Println(pwettyPwint("IP: "+remoteAddr.String(), textProperties{Color: "#eba434", Bold: true}))
	if err != nil {
		fmt.Println(err)
	}
	return
}

var SERVER_ADDR_IPV4 *net.UDPAddr = &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 53}
var SERVER_ADDR_IPV6 *net.UDPAddr = &net.UDPAddr{IP: net.ParseIP("localhost"), Port: 53}

var ONEDOT_ADDR_IPV4 *net.UDPAddr = &net.UDPAddr{IP: net.ParseIP("1.1.1.1"), Port: 53}
var ONEDOT_ADDR_IPV6 *net.UDPAddr = &net.UDPAddr{IP: net.ParseIP("2606:4700:4700::1111"), Port: 53}

var blocklist = new(StringSet)

func main() {
	blockMode := false
	for i, v := range os.Args {
		if v == "-m" && i+1 < len(os.Args) { //op mode setting
			if os.Args[i+1] == "block" {
				blockMode = true
			} else if os.Args[i+1] != "forward" {
				fmt.Println("Invalid mode")
				os.Exit(0)
			}
		}
	}
	App := webApp{
		ToClient: make(chan string, 30),
	}

	connection_IPV4, err_IPV4 := net.ListenUDP("udp", SERVER_ADDR_IPV4)
	connection_IPV6, err_IPV6 := net.ListenUDP("udp", SERVER_ADDR_IPV6)
	/* App.connection_IPV4 = connection_IPV4
	App.connection_IPV6 = connection_IPV6 */
	fmt.Println(connection_IPV4.LocalAddr(), connection_IPV6.LocalAddr()) // current
	if err_IPV4 != nil {
		fmt.Println(err_IPV4)
	}
	if err_IPV6 != nil {
		fmt.Println(err_IPV6)
	}
	if blockMode {
		blocklist.loadDataFromFile("Blocklist Management/blocklist.txt")
		go dualStart_blocking(connection_IPV4, connection_IPV6, &App)
	} else {
		go dualStart_forwarding(connection_IPV4, connection_IPV6, &App)
	}
	StartWAServer(&App)
}
