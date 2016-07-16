// +build linux darwin

package monitor_agent

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func checkProcess(service *SvcConfig) {
	proc_regexp := "ps -ef | grep " + service.Parameters["proc_regexp"] + " | wc -l"
	linecnt, err := strconv.Atoi(service.Parameters["line_count"])
	if err != nil {
		panic(err)
	}
	if len(proc_regexp) == 0 || linecnt < 1 {
		service.Status = "Error"
		fmt.Printf("Missing proc_regexp or line count for %s\n", service.Name)
		return
	}
	cmd := exec.Command("/bin/bash", "-c", proc_regexp)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
//	fmt.Printf("Command returned %s (%d bytes)\n", out.String(), len(out.String()))
	nl, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		panic(err)
	}
	if nl < linecnt {
		service.Status = "Not found"
	} else {
		service.Status = "OK"
	}
}

func checkFile(service *SvcConfig) {
	filename := service.Parameters["filename"]
	younger_than_hours, err := strconv.ParseFloat(service.Parameters["younger_than_hours"], 64)
	if err != nil {
		panic(err)
	}
	if len(filename) == 0 || younger_than_hours < 0.001 {
		service.Status = "Error"
		fmt.Printf("Missing filename or age limit for %s\n", service.Name)
		return
	}
	stat, err := os.Stat(filename)
	if err != nil {
		service.Status = "Not found"
	} else {
		if time.Since(stat.ModTime()).Hours() > younger_than_hours {
			service.Status = "Old"
		} else {
			service.Status = "OK"
		}
	}
}
