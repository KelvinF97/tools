package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	//源端口，目的端口
	var fromPort, toPort int = 5389, 22
	fromAddr := fmt.Sprintf("0.0.0.0:%d", fromPort)
	toAddr := fmt.Sprintf("127.0.0.1:%d", toPort)

	fromListener, err := net.Listen("tcp", fromAddr)

	if err != nil {
		log.Fatalf("Unable to listen on: %s, error: %s\n", fromAddr, err.Error())
	}

	// if you use defer fromListener.Close(), it's would warn you modify this method
	// otherwise you could write like the bottom
	defer func(fromListener net.Listener) {
		_ = fromListener.Close()
	}(fromListener)

	for {
		fromCon, err := fromListener.Accept()
		if err != nil {
			fmt.Printf("Unable to accept a request, error: %s\n", err.Error())
		} else {
			fmt.Println("new connect:" + fromCon.RemoteAddr().String())
		}

		//这边最好也做个协程，防止阻塞
		toCon, err := net.Dial("tcp", toAddr)
		if err != nil {
			fmt.Printf("can not connect to %s\n", toAddr)
			continue
		}
		go handleConnection(fromCon, toCon)
		go handleConnection(toCon, fromCon)
	}
}

func handleConnection(r, w net.Conn) {
	defer func(r net.Conn) {
		err := r.Close()
		if err != nil {
			fmt.Println("read operation had closed.")
		}
	}(r)
	defer func(w net.Conn) {
		err := w.Close()
		if err != nil {
			fmt.Println("write operation had closed.")
		}
	}(w)

	var buffer = make([]byte, 100000)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			break
		}

		n, err = w.Write(buffer[:n])
		if err != nil {
			break
		}
	}

}


