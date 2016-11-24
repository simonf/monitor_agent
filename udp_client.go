package monitor_agent

import (
	"fmt"
	"net"
)

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
		_, addr, err := socket.ReadFromUDP(buf)
		// fmt.Println("Received ", string(buf[0:n]), " from ", addr)

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
	if server_config != nil {
		server_list = make([]net.IP, 0)
		ips_from_dns, err := net.LookupIP(server_config.Name)
		if err != nil {
			server_list = append(server_list, net.ParseIP(server_config.Address))
		} else {
			server_list = append(server_list, ips_from_dns[0])
		}
	}
	sl_mutex.Lock()
	for _, ip := range server_list {
		ua := net.UDPAddr{IP: ip, Port: 41237}
		// fmt.Printf("Dialing UDPAddress %v\n", ip)
		conn, err := net.DialUDP("udp4", nil, &ua)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// fmt.Printf("Sending %d bytes\n", len(b))

		_, err = conn.Write(b)

		if err != nil {
			fmt.Println(err)
		}
		conn.Close()
	}
	sl_mutex.Unlock()
	// fmt.Println("Send complete")
}
