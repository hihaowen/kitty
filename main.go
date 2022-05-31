package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	service := ":16106"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		handleClient(conn)

		conn.Close() // we're finished with this client
	}
}

func handleClient(conn net.Conn) {
	var buf [1]byte
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			return
		}

		reqBytes := buf[:]
		reqString := string(reqBytes)

		log.Println(n, buf[:], reqString)

		if bytes.Compare(buf[:], []byte("quit")) == 0 {
			log.Println(buf)
			return
		}

		_, err2 := conn.Write(reqBytes)
		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
