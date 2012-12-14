package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

var (
	addr           = flag.String("addr", "0.0.0.0:843", "address")
	file           = flag.String("file", "", "the socket policy file")
	request        = []byte("<policy-file-request/>\x00")
	response       = []byte("<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\" /></cross-domain-policy>\x00")
	buf            = make([]byte, len(request))
	readWriteLimit = 5 * time.Second
)

func main() {
	flag.Parse()

	var err error

	if *file != "" {
		response, err = ioutil.ReadFile(*file)
		if err != nil {
			log.Fatal(err)
		}
		response = append(response, 0)
	}

	run(make(chan bool, 1))
}

func run(startedChan chan bool) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	listner, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	startedChan <- true
	log.Print("=== Flash Socket Policy File Server ===")

	for {
		conn, err := listner.AcceptTCP()
		if err == nil {
			go loop(conn)
		}
	}
}

func loop(conn *net.TCPConn) {
	defer conn.Close()

	conn.SetLinger(5)
	conn.SetKeepAlive(false)
	conn.SetNoDelay(true)
	now := time.Now()

	conn.SetReadDeadline(now.Add(readWriteLimit))

	if _, err := io.ReadFull(conn, buf); err == nil {
		conn.Write(response)
	}
}
