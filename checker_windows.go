// +build windows
package monitor_agent

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

import ps "github.com/mitchellh/go-ps"

func checkProcess(service *SvcConfig) {
	proc_regexp := service.Parameters["proc_regexp"]
	linecnt, err := strconv.Atoi(service.Parameters["line_count"])
	if err != nil {
		panic(err)
	}
	if len(proc_regexp) == 0 || linecnt < 1 {
		service.Status = "Error"
		fmt.Printf("Missing proc_regexp or line count for %s\n", service.Name)
		return
	}
	plist, err := ps.Processes()
	if err != nil {
		panic(err)
	}
	cnt := 0
	for _,proc := range plist {
		tf, err := regexp.MatchString(proc_regexp,proc.Executable())
		if err != nil {
			panic(err)
		}
		if tf {
			cnt += 1
		}
	}

	if cnt < linecnt {
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