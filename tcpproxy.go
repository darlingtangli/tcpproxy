package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func handle(localConn net.Conn, remoteConn net.Conn) {
	var once sync.Once
	onceBody := func() {
		localConn.Close()
		remoteConn.Close()
	}
	done := make(chan struct{})

	// handle upstream
	go func() {
		defer close(done)
		nUpsteam, _ := io.Copy(localConn, remoteConn)
		once.Do(onceBody)
		log.Printf("upstream %d bytes: %v <--  %v\n", nUpsteam, localConn.RemoteAddr(), remoteConn.RemoteAddr())
	}()

	// handle downstream
	nDownstream, _ := io.Copy(remoteConn, localConn)
	once.Do(onceBody)
	log.Printf("downstream %d bytes: %v -->  %v\n", nDownstream, localConn.RemoteAddr(), remoteConn.RemoteAddr())

	<-done

	return
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("usage: %s <local port> <remote host> <remote port>\n", os.Args[0])
	}

	localAddr := ":" + os.Args[1]
	remoteAddr := os.Args[2] + ":" + os.Args[3]

	l, err := net.Listen("tcp4", localAddr)
	if err != nil {
		log.Fatalf("bind %s fail: %s\n", localAddr, err.Error())
	}

	for {
		localConn, err := l.Accept()
		if err != nil {
			log.Printf("accept fail: %s\n", err.Error())
			continue
		}
		log.Println("accept:", localConn.RemoteAddr())

		remoteConn, err := net.Dial("tcp4", remoteAddr)
		if err != nil {
			log.Printf("dail remote addr %s fail: %s", remoteAddr, err.Error())
			localConn.Close()
			continue
		}
		log.Println("connect remote:", remoteConn.RemoteAddr())

		go handle(localConn, remoteConn)
	}

	return
}
