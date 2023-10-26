package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const quote = `
The most important thing in the programming language is the name. 
A language will not succeed without a good name. I have recently invented a very
good name and now I am looking for a suitable language.
`

var small []byte // ~2k chunk
var large []byte // ~40k chunk

func printUsage() {
	prog := "duplex-stress"
	if len (os.Args) > 0 {
		prog = os.Args[0]
	}
	fmt.Printf("Usage: %s(client|server) ipv4addr port\n\n", prog)
	fmt.Println("Examples:")
	fmt.Printf("%s server 0.0.0.0 9191\n\n", prog)
	fmt.Printf("%s client 127.0.0.1 9191\n", prog)
}

func setup() {
	small = make([]byte, 10 * len(quote))
	for i := 0; i < 10; i++ {
		copy(small[(i * len(quote)):], []byte(quote))
	}

	large = make([]byte, 20 * len(small))
	for i := 0; i < 20; i++ {
		copy(large[i*len(small):], small)
	}
}

func main() {

	if len(os.Args) < 4 {
		printUsage()
		os.Exit(1)
	}

	setup()


	switch os.Args[1] {
	case "client":
		startClient(os.Args[2], os.Args[3])
	case "server":
		startServer(os.Args[2], os.Args[3])	
	}

	for {
		time.Sleep(1000 * time.Second)
	}
}


func receiver(conn net.Conn, ready chan bool, signal chan bool) {
	buf := make([]byte, len(small))

	for {
		// read 1 for every 20 (~ 2k blocks)
		for i := 0; i < 20; i ++ {
			<- ready
			_, err := io.ReadFull(conn, buf)
			if err != nil {
				panic(err)
			}
			fmt.Print("r");
			signal <- true
		}

		// catchup reading 20 * 20 (reading ~40k in 2k chunks)
		<- ready
		for i := 0; i < 400; i++ {
			_, err := io.ReadFull(conn, buf)
			if err != nil {
				panic(err)
			}
			fmt.Print("r");
		}
		signal <- true
	}
}


func sender(conn net.Conn, ready chan bool, signal chan bool) {
	for {
		<- ready
		// write the equivalent of 20 entries in one ~40k chunk
		_, err := conn.Write(large)
		if err != nil {
			panic(err)
		}

		fmt.Print("w")
		signal <- true
	}
}

func startClient(ip string, port string) {
	conn, err := net.Dial("tcp4", ip + ":" + port)
	if err != nil {
		panic(err)
	}

	rcv := make(chan bool)
	snd := make(chan bool)

	go receiver(conn, rcv, snd)
	go sender(conn, snd, rcv)

	// sender goes first
	snd <- true
}

func startServer(ip string, port string) {
	listener, err := net.Listen("tcp4", ip + ":" + port)
	if err != nil {
		panic(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	rcv := make(chan bool)
	snd := make(chan bool)

	go receiver(conn, rcv, snd)
	go sender(conn, snd, rcv)

	// receiver goes first
	rcv <- true
}