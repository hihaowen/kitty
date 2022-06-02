package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type getMap struct {
	val map[string]string
}

var getMaps getMap

func init() {
	getMaps = *new(getMap)
	getMaps.val = map[string]string{}
}

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
	for {
		var buf [512]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			return
		}

		reqBytes := buf[:]
		reqBytes = bytes.TrimRight(reqBytes, "\x00\x0a\x0d")
		reqString := string(reqBytes)

		log.Printf("用户输入 %d 个字节: %s", n, reqString)
		log.Println(reqBytes)

		if reqString == "quit" {
			log.Println("退出咯...")
			return
		}

		args := strings.Split(reqString, " ")
		if len(args) < 2 {
			conn.Write([]byte("\n"))
			continue
		}

		cmd := args[0]
		key := args[1]

		log.Println(args, len(args))

		switch {
		case strings.ToUpper(cmd) == "GET":
			val, isSet := getMaps.val[key]
			if ! isSet {
				conn.Write([]byte("key is not set: " + key + "\n"))
				break
			}
			conn.Write([]byte(val + "\n"))
			break
		case strings.ToUpper(cmd) == "SET":
			if len(args) < 3 {
				conn.Write([]byte("value is not set: " + key + "\n"))
				break
			}
			getMaps.val[key] = args[2]
			conn.Write([]byte("ok\n"))
			break
		case strings.ToUpper(cmd) == "DEL":
			_, isSet := getMaps.val[key]
			if ! isSet {
				conn.Write([]byte("key is not set: " + key + "\n"))
				break
			}
			delete(getMaps.val, key)
			conn.Write([]byte("ok\n"))
			break
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
