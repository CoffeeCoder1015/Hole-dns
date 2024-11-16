package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	dns "github.com/miekg/dns"
)

func CreateNullReponse(originalMsg dns.Msg) dns.Msg {
	Header := dns.RR_Header{
		Name:   originalMsg.Question[0].Name,
		Rrtype: originalMsg.Question[0].Qtype,
		Class:  dns.ClassINET,
		Ttl:    3600,
	}
	if originalMsg.Question[0].Qtype == 1 { //ipv4
		nilAnswer := new(dns.A)
		nilAnswer.Hdr = Header
		nilAnswer.A = net.ParseIP("0.0.0.0")
		originalMsg.Answer = append(originalMsg.Answer, nilAnswer)
		return originalMsg
	} else if originalMsg.Question[0].Qtype == 28 { //ipv6
		nilAnswer := new(dns.AAAA)
		nilAnswer.Hdr = Header
		nilAnswer.AAAA = net.ParseIP("0.0.0.0")
		originalMsg.Answer = append(originalMsg.Answer, nilAnswer)
		return originalMsg
	}
	return dns.Msg{}
}

type null interface{}

type StringSet struct {
	data map[string]null //the struct is null essentially
}

func (s *StringSet) loadDataFromFile(fileName string) {
	s.data = make(map[string]null)
	file, _ := os.Open(fileName)
	fileReader := bufio.NewReader(file)
	for {
		str, err := fileReader.ReadString('\n')
		str = strings.TrimRight(str, "\r\n ")
		s.data[str] = struct{}{}
		if err == io.EOF {
			break
		}
	}
	fmt.Println(pwettyPwint("Hole DNS pre-loader:", textProperties{Bold: true, Color: "#32a8a8"}), len(s.data), "urls loaded")
}

func (s *StringSet) isIn(str string) bool {
	_, ok := s.data[str]
	return ok
}

//for checking if url is blocked:
//best practical solution (in order) : binary search, interpolation search
//Interesting and viable alternative solutions: Bloom filter, Merkle tree
