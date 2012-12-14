package main

import (
	"io/ioutil"
	"net"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	startedChan := make(chan bool, 1)
	go func() {
		run(startedChan)
	}()

	<-startedChan

	clientCount := 2000

	resultChan := make(chan interface{}, clientCount)
	b.StartTimer()
	for i := 0; i < clientCount; i++ {
		go doGet(resultChan)
	}

	for i := 0; i < clientCount; i++ {
		<-resultChan
	}

	b.StopTimer()
}

func doGet(result chan interface{}) {
	conn, err := net.Dial("tcp", "127.0.0.1:843")

	if err != nil {
		result <- err
		return
	}

	defer conn.Close()
	conn.Write(request)
	_, err = ioutil.ReadAll(conn)
	if err != nil {
		result <- err
	} else {
		result <- nil
	}
}
