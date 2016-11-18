package monitor_agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ServerConfig struct {
	Name    string
	Address string
}

func (s ServerConfig) print() {
	fmt.Printf("Name: %s, Address: %s\n", s.Name, s.Address)
}

func readServer(filename string) (*ServerConfig, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	tmp := ServerConfig{}

	err = json.Unmarshal(b, &tmp)
	if err != nil {
		return nil, err
	}

	return &tmp, nil
}
