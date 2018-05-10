package main

import (
	"fmt"
	"net"
	"os"
	sig "os/signal"
)

func main() {
	addr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 3000,
		Zone: "::1",
	}

	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listened on", addr.IP, "on port", addr.Port)
	/*
		catch interrupt for SIGINT (Ctrl + C) and immediately exit
		as on Windows, when it's exit without handled will give
		exit status 2
	*/
	interrupt := make(chan os.Signal, 1)
	sig.Notify(interrupt, os.Interrupt)
	running := true
	go func() {
		<-interrupt
		running = false
		os.Exit(0)
	}()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		go func(conn *net.TCPConn) {
			fmt.Println("accepted connection")
			buffer := make([]byte, 512)
			bytesread, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Read connection error", err)
			} else {
				fmt.Printf("Getting %d chars read\n", bytesread)
				fmt.Printf("%s", buffer[:bytesread])
				fmt.Println("echoing back")
				_, _ = conn.Write([]byte("echo> "))
				_, _ = conn.Write(buffer[:bytesread])
			}
			conn.Close()

		}(conn)
		if !running {
			break
		}
	}
}
