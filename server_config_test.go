package monitor_agent

import (
	"testing"
)

func TestLoadsServerFile(t *testing.T) {
	sc, err := readServer("settings.json")
	if err != nil {
		t.Error(err)
	}
	if sc.Address != "123.145.167.189" {
		t.Error("Address not equal")
	}
}
