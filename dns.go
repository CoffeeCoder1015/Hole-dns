// manually implemented dns parsing
// i cant get self referencing to work so i ditched it :C
package main

import (
	"fmt"
	"math"
	"net"
	"net/url"
	"strconv"

	dns "github.com/miekg/dns"
)

var TRUE string = pwettyPwint("True", textProperties{Color: "#03fc62"})
var FALSE string = pwettyPwint("False", textProperties{Color: "#fc3903"})

func getUrlIP(URL string) net.IP {
	parsedURL, _ := url.Parse(URL)
	if parsedURL.Host[len(parsedURL.Host)-1] != '.' {
		parsedURL.Host += "."
	}

	dnsMessage := dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:               dns.Id(),
			Response:         false,
			Opcode:           dns.OpcodeQuery,
			RecursionDesired: true,
			Rcode:            0,
		},
		Question: make([]dns.Question, 1),
	}
	dnsMessage.Question[0] = dns.Question{
		Name:   parsedURL.Host,
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}

	data, err := dnsMessage.Pack()
	fmt.Println(data)
	if err != nil {
		fmt.Println("pack err:", err)
		return nil
	}

	fConn, err := net.Dial("udp", ONEDOT_ADDR_IPV4.String())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fConn.Write(data)
	output := make([]byte, 512)
	fConn.Read(output)

	unpackErr := dnsMessage.Unpack(output)
	if unpackErr != nil {
		fmt.Println(unpackErr)
		return nil
	}
	fmt.Println(dnsMessage)
	for _, a := range dnsMessage.Answer {
		if a.Header().Rrtype == dns.TypeA {
			IPV4_record := a.(*dns.A)
			return IPV4_record.A
		}
	}
	return nil
}

func bin(n byte) string {
	return strconv.FormatInt(int64(n), 2)
}

// move this into snippets
func readBits(data []byte, at_index int, n int) int64 {
	arr_start_index := int(math.Floor(float64(at_index) / 8))
	rel_index := at_index - 8*arr_start_index

	rNum := int64(0)
	for i := arr_start_index; i < len(data); i++ {
		B := data[i]
		if (rel_index + n) >= 8 {
			mask_length := 8 - rel_index
			nB := B & ((1 << mask_length) - 1)
			rNum += int64(nB) << (n - mask_length)
			n -= mask_length
			rel_index = 0
		} else {
			//mask_length := n
			bit_shift := 8 - (n + rel_index)
			nB := (B >> byte(bit_shift)) & ((1 << n) - 1)
			rNum += int64(nB)
			break
		}
	}
	return rNum
}

func boolToStr(b bool) string {
	if b {
		return TRUE
	}
	return FALSE
}

func formatDomainName(s string) string {
	return pwettyPwint(s, textProperties{Italics: true, Color: "#5bf0a5"})
}

func formatTypeNames(s string) string {
	return pwettyPwint(s, textProperties{Color: "#f73e94"})
}

func displayMessage(parsedMsg dns.Msg) {
	fmt.Println(pwettyPwint("-----  Header  -----", textProperties{Bold: true, Color: "#168bd9"}))
	fmt.Println("ID:", parsedMsg.Id)
	fmt.Println("IsResponse:", boolToStr(parsedMsg.Response))
	fmt.Println("OpCode:", formatTypeNames(dns.OpcodeToString[parsedMsg.Opcode]), "Authoritative:", boolToStr(parsedMsg.Authoritative), "Truncated:", boolToStr(parsedMsg.Truncated))
	fmt.Println("Recursion Desired:", boolToStr(parsedMsg.RecursionDesired), "RecursionAvailable:", boolToStr(parsedMsg.RecursionAvailable))
	fmt.Println("RCode:", parsedMsg.Rcode)
	if len(parsedMsg.Question) > 0 {
		Q := parsedMsg.Question[0]
		fmt.Println(pwettyPwint("	-----  Question  -----", textProperties{Bold: true, Color: "#16d977"}))
		fmt.Println(" 	", formatDomainName(Q.Name), "Qtype:", formatTypeNames(dns.TypeToString[Q.Qtype]), "QClass:", dns.ClassToString[Q.Qclass])
	}
	fmt.Println()
	for a := 0; a < len(parsedMsg.Answer); a++ {
		fmt.Println(pwettyPwint(fmt.Sprintf("	-----  Answer(%d)  -----", a), textProperties{Bold: true, Color: "#7e18cc"}))
		A := parsedMsg.Answer[a]
		A_Hddr := A.Header()
		fmt.Println(" 	", formatDomainName(A_Hddr.Name))
		fmt.Println("  	RR Class:", dns.ClassToString[A_Hddr.Class],
			"RR Type:", formatTypeNames(dns.TypeToString[A_Hddr.Rrtype]),
			pwettyPwint("TTL:", textProperties{Color: "#aaf542"}), A_Hddr.Ttl)
	}
	fmt.Println()
}
