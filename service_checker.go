package monitor_agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SvcConfig struct {
	Name         string
	Status       string
	Service_type string
	Parameters   map[string]string
}

func (s SvcConfig) print() {
	fmt.Printf("Name: %s, Status: %s, Service_type: %s Parameters:\n", s.Name, s.Status, s.Service_type)
	for k, v := range s.Parameters {
		fmt.Printf("\t%s: %s\n", k, v)
	}
}

func readServices(filename string) []*SvcConfig {
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	tmp := struct {
		Services []SvcConfig
	}{Services: make([]SvcConfig, 0)}

	err = json.Unmarshal(b, &tmp)
	if err != nil {
		panic(err)
	}
	retval := make([]*SvcConfig, 0)
	for i, _ := range tmp.Services {
		psvc := &tmp.Services[i]
		retval = append(retval, psvc)
	}
	return retval
}

func checkServices(services []*SvcConfig) {
	for _, psvc := range services {
		// fmt.Printf("Checking %s\n", psvc.Name)
		switch psvc.Service_type {
		case "process":
			checkProcess(psvc)
		case "file":
			checkFile(psvc)
		default:
			fmt.Println("No Service_type match")
		}
	}
}

