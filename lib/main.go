package lib

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const port = "8080"

// RunHost takes an IP address as argument
// and listens for connection on that IP
func RunHost(ip string) {

	var wg sync.WaitGroup
	wg.Add(2)

	ipAndPort := ip + ":" + port
	listener, listenErr := net.Listen("tcp", ipAndPort)
	if listenErr != nil {
		log.Fatal("Error: ", listenErr)
	}
	fmt.Println("Listening on", ipAndPort)

	conn, acceptErr := listener.Accept()
	if acceptErr != nil {
		log.Fatal("Error: ", acceptErr)
	}
	fmt.Println("New Connection Accepted")

	go handleHostRead(conn, &wg)
	go handleHostSend(conn, &wg)

	wg.Wait()

}

func handleHostRead(conn net.Conn, wg *sync.WaitGroup) {
	for {
		reader := bufio.NewReader(conn)
		message, readErr := reader.ReadString('\n')
		if readErr != nil {
			log.Fatal("Error", readErr)
			wg.Done()
		}
		fmt.Println("Message received: ", message)
	}

}

func handleHostSend(conn net.Conn, wg *sync.WaitGroup) {
	for {
		replyReader := bufio.NewReader(os.Stdin)
		replyMessage, replyErr := replyReader.ReadString('\n')
		if replyErr != nil {
			log.Fatal("Error", replyErr)
			wg.Done()
		}
		fmt.Fprint(conn, replyMessage)
	}
}

// RunGuest takes a destination IP as an
// argument and connects to that IP
func RunGuest(ip string) {
	var wg sync.WaitGroup
	wg.Add(2)

	ipAndPort := ip + ":" + port
	conn, dialErr := net.Dial("tcp", ipAndPort)
	if dialErr != nil {
		log.Fatal("Error: ", dialErr)
	}
	go handleGuestRead(conn, &wg)
	go handleGuestSend(conn, &wg)

	wg.Wait()

}

func handleGuestRead(conn net.Conn, wg *sync.WaitGroup) {
	for {
		replyReader := bufio.NewReader(conn)
		replyMessage, replyErr := replyReader.ReadString('\n')
		if replyErr != nil {
			log.Fatal("Error: ", replyErr)
			wg.Done()
		}
		fmt.Println("Message received:", replyMessage)
	}
}

func handleGuestSend(conn net.Conn, wg *sync.WaitGroup) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, readErr := reader.ReadString('\n')
		if readErr != nil {
			log.Fatal("Error: ", readErr)
			wg.Done()
		}
		fmt.Fprint(conn, message)
	}
}
