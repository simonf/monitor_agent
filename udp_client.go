package monitor_agent

import (
	"fmt"
	"net"
	"sync"
)

var server_list = make([]net.IP, 0)
var sl_mutex = &sync.Mutex{}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
}

func listenForServers() {
	addr, err := net.ResolveUDPAddr("udp4", ":41238")
	checkError(err)

	socket, err := net.ListenUDP("udp4", addr)
	// socket, err := net.ListenUDP("udp4", &net.UDPAddr{
	//   IP:   net.IPv4(0, 0, 0, 0),
	//   Port: 41238,
	// })

	defer socket.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := socket.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			match := false
			sl_mutex.Lock()
			for _, a := range server_list {
				if a.String() == addr.IP.String() {
					match = true
				}
			}
			if !match {
				server_list = append(server_list, addr.IP)
			}
			sl_mutex.Unlock()
		}
	}
}

func sendToServers(b []byte) {
	sl_mutex.Lock()
	for _, ip := range server_list {
		ua := net.UDPAddr{IP: ip, Port: 41237}
		fmt.Printf("Dialing UDPAddress %v", ip)
		conn, err := net.DialUDP("udp4", nil, &ua)
		if err != nil {
			panic(err)
		}
		_, err = conn.Write(b)

		if err != nil {
			fmt.Println(err)
		}
		conn.Close()
	}
	sl_mutex.Unlock()
}
