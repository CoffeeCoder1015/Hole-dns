package main

import (
	"fmt"
	"net"
	"time"

	dns "github.com/miekg/dns"
)

func dualStart_blocking(connIPV4 *net.UDPConn, connIPV6 *net.UDPConn, sendChan *webApp) {
	sendChan.blockOn = true
	sendChan.SendToClient(fmt.Sprintf("b%t", sendChan.blockOn))
	go connectionHandler_blocking(connIPV4, ONEDOT_ADDR_IPV4.String(), sendChan)
	go connectionHandler_blocking(connIPV6, ONEDOT_ADDR_IPV6.String(), sendChan)

}

func dualStart_forwarding(connIPV4 *net.UDPConn, connIPV6 *net.UDPConn, sendChan *webApp) {
	sendChan.SendToClient(fmt.Sprintf("b%t", sendChan.blockOn))
	go connectionHandler_forwarding(connIPV4, ONEDOT_ADDR_IPV4.String(), sendChan)
	go connectionHandler_forwarding(connIPV6, ONEDOT_ADDR_IPV6.String(), sendChan)
}

func connectionHandler_blocking(connection *net.UDPConn, forwardingAddress string, sendChan *webApp) {
	for {
		data, Addr := readUDP(connection, 512)
		go func() {
			sendChan.SendToClient(connection.LocalAddr().String())

			parsedMsg := dns.Msg{}
			parsedMsg.Unpack(data)
			//checks if url is in blocklist
			url := parsedMsg.Question[0].Name
			if url[len(url)-1] == '.' {
				url = url[:len(url)-1]
			}
			sendChan.SendToClient(fmt.Sprintf("co %s", url))

			if blocklist.isIn(url) {
				parsedMsg = CreateNullReponse(parsedMsg)
				data, _ = parsedMsg.Pack()
				connection.WriteToUDP(data, Addr)
				fmt.Println(url, "is blocked")
				return
			}

			//1.1.1.1 forwarding
			fConn, err := net.Dial("udp", forwardingAddress)
			if err != nil {
				sendChan.SendToClient(err.Error())
				return
			}
			fConn.SetReadDeadline(time.Now().Add(1 * time.Second))
			fConn.Write(data)
			fConn.Read(data)

			/* unpack msg from cloudflare here parsedMsg.Unpack(data)*/
			//resp to querier
			connection.WriteToUDP(data, Addr)
		}()
	}
}

func connectionHandler_forwarding(connection *net.UDPConn, forwardingAddress string, sendChan *webApp) {
	for {
		data, Addr := readUDP(connection, 512)
		go func() {
			sendChan.SendToClient(connection.LocalAddr().String())
			parsedMsg := dns.Msg{}
			parsedMsg.Unpack(data)

			url := parsedMsg.Question[0].Name
			if url[len(url)-1] == '.' {
				url = url[:len(url)-1]
			}
			sendChan.SendToClient(fmt.Sprintf("co %s", url))

			//1.1.1.1 forwarding
			fConn, err := net.Dial("udp", forwardingAddress)
			if err != nil {
				sendChan.SendToClient(err.Error())
				return
			}
			fConn.SetReadDeadline(time.Now().Add(1 * time.Second))
			fConn.Write(data)
			fConn.Read(data)

			/* unpack msg from cloudflare here parsedMsg.Unpack(data)*/
			//resp to querier
			connection.WriteToUDP(data, Addr)
		}()
	}
}
