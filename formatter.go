package monitor_agent

import (
	"fmt"
	"os"
	"simonf.net/monitor_db"
	"time"
)

func makeAgentPayload(c *monitor_db.Computer) []byte {
	return []byte(c.JSON())
}

func servicesFromSvcConfigs(svcs []*SvcConfig) []*monitor_db.Service {
	retval := make([]*monitor_db.Service, 0)
	for _, sc := range svcs {
		retval = append(retval, serviceFromConfig(sc))
	}
	return retval
}

func serviceFromConfig(sc *SvcConfig) *monitor_db.Service {
	return &monitor_db.Service{Name: sc.Name, Status: sc.Status, Updated: time.Now()}
}

func computerFromServices(svcs []*monitor_db.Service) *monitor_db.Computer {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Unable to get hostname")
		hostname = "Unknown"
	}
	return &monitor_db.Computer{Name: hostname, Status: calcHostStatus(svcs), Services: svcs, Updated: time.Now()}
}

func calcHostStatus(svcs []*monitor_db.Service) string {
	any_ok := false
	any_bad := false
	for _, svc := range svcs {
		if svc.Status == "OK" {
			any_ok = true
		}
		if svc.Status != "OK" {
			any_bad = true
		}
	}
	if any_ok {
		if !any_bad {
			return "Green"
		} else {
			return "Amber"
		}
	} else {
		return "Red"
	}
}
