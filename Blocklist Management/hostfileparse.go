package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
)

func lineParser(s string, outFile *os.File) {
	if len(s) <= 8 || s[0] == '#' {
		return
	}
	if s[:7] == "0.0.0.0" {
		ip := s[8:]
		outFile.WriteString(ip)
	}
}

func main() {
	hostfile, _ := os.Open("default_host_raw.txt")
	defer hostfile.Close()
	outfile, _ := os.OpenFile("blocklist.txt", os.O_WRONLY, fs.ModeAppend)
	defer outfile.Close()

	FileReader := bufio.NewReader(hostfile)

	var wg sync.WaitGroup
	mutex := sync.Mutex{}
	for {
		s, err := FileReader.ReadString('\n')
		wg.Add(1)

		go func() {
			defer mutex.Unlock()
			defer wg.Done()
			mutex.Lock()
			lineParser(s, outfile)
		}()

		if err == io.EOF {
			fmt.Println(err)
			break
		}
	}
	wg.Wait()
}
