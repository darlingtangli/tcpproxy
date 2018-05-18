/**
* @file tcpproxy.go
* @brief
* @author darlintangli@126.com
* @version 1.0
* @date 2018-05-18
 */
package main

import (
	"log"
	"net"
	"os"
)

func handle(c net.Conn, remoteAddr string) {
	remoteConn, err := net.Dial("tcp4", remoteAddr)
	if err != nil {
		log.Fatalf("dail remote addr %s fail: %s", remoteAddr, err.Error())
	}
	log.Println("connect remote:", remoteConn.RemoteAddr())

	writeHandlerClosed := make(chan struct{})
	// handle write
	go func() {
		defer close(writeHandlerClosed)
		remoteReadBuf := make([]byte, 10240)
		for {
			nUpsteam, err := remoteConn.Read(remoteReadBuf)
			if err != nil {
				log.Printf("read remote fail: %s\n", err.Error())
				break
			}
			nUpsteam, err = c.Write(remoteReadBuf[:nUpsteam])
			if err != nil {
				log.Printf("upstream fail: %s\n", err.Error())
				break
			}
			log.Printf("upstream %d bytes <-- remoteAddr %s\n", nUpsteam, remoteAddr)
		}
		log.Println("write handler exit...")
	}()

	// handle read
	localReadBuf := make([]byte, 10240)
	for {
		nDownstream, err := c.Read(localReadBuf)
		if err != nil {
			log.Printf("read local fail: %s\n", err.Error())
			break
		}
		nDownstream, err = remoteConn.Write(localReadBuf[:nDownstream])
		if err != nil {
			log.Printf("downstream fail: %s\n", err.Error())
			break
		}
		log.Printf("downstream %d bytes --> remote addr %s\n", nDownstream, remoteAddr)
	}
	log.Println("read handler exit...")
	c.Close()
	remoteConn.Close()
	<-writeHandlerClosed

	return
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("usage: %s <local port> <remote host> <remote port>\n", os.Args[0])
	}

	localAddr := ":" + os.Args[1]
	l, err := net.Listen("tcp4", localAddr)
	if err != nil {
		log.Fatalf("bind %s fail: %s\n", localAddr, err.Error())
	}

	remoteAddr := os.Args[2] + ":" + os.Args[3]
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("accept fail: %s\n", err.Error())
		}
		log.Println("accept:", c.RemoteAddr())
		go handle(c, remoteAddr)
	}

	return
}
