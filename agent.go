package monitor_agent

import (
	"fmt"
	"time"
)

var config_file = "services.json"
var sleep_minutes = 1

func Start() {
	go listenForServers()
	go periodicallyCheckServices()
}

func periodicallyCheckServices() {
	for {
		svcs := readServices(config_file)
		if len(svcs) < 1 {
			panic("No services to monitor")
		}
		checkServices(svcs)
		fmt.Printf("Checked %d services\n", len(svcs))
		svcstosend := servicesFromSvcConfigs(svcs)
		c := computerFromServices(svcstosend)
		fmt.Printf("Sending to %d servers\n", len(server_list))
		ba := makeAgentPayload(c)
		sendToServers(ba)
		time.Sleep(time.Duration(sleep_minutes) * time.Minute)
		// printAll(svcs)
	}
}

func printAll(svcs []*SvcConfig) {
	for _, svc := range svcs {
		svc.print()
	}
}
