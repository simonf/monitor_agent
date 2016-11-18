package monitor_agent

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var server_file = "settings.json"
var config_file = "services.json"
var sleep_minutes = 1

var server_list = make([]net.IP, 0)
var sl_mutex = &sync.Mutex{}

func Start() {
	sc, err := readServer(server_file)
	if err != nil {
		go listenForServers()
	} else {
		server_list = append(server_list, net.ParseIP(sc.Address))
	}
	go periodicallyCheckServices()
}

func periodicallyCheckServices() {
	for {
		svcs := readServices(config_file)
		if len(svcs) < 1 {
			panic("No services to monitor")
		}
		checkServices(svcs)
		// fmt.Printf("Checked %d services\n", len(svcs))
		svcstosend := servicesFromSvcConfigs(svcs)
		c := computerFromServices(svcstosend)
		// fmt.Printf("Sending to %d servers\n", len(server_list))
		ba := makeAgentPayload(c)
		sendToServers(ba)
		time.Sleep(time.Duration(sleep_minutes) * time.Minute)
		// printAll(svcs)
	}
}

func RunOnce() {
	svcs := readServices(config_file)
	if len(svcs) < 1 {
		panic("No services to monitor")
	}
	checkServices(svcs)
	fmt.Printf("Checked %d services\n", len(svcs))
	svcstosend := servicesFromSvcConfigs(svcs)
	c := computerFromServices(svcstosend)
	ba := makeAgentPayload(c)
	fmt.Printf("%s\n", string(ba))
}

func printAll(svcs []*SvcConfig) {
	for _, svc := range svcs {
		svc.print()
	}
}
