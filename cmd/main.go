package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/niubir/sysproxy"
)

const version = "1.0.8"

var (
	on         string
	off        string
	setIgnore  bool
	listenQuit bool

	serviceRegex = regexp.MustCompile(`(\w+)=([\d\.]+):(\d+)`)
)

type Service struct {
	Original string
	Protocol string
	Addr     sysproxy.Addr
}

func init() {
	flag.Usage = func() {
		fmt.Println("sysproxy version", version)
		flag.PrintDefaults()
	}
	flag.StringVar(&on, "on", "", "Turn on services in the format: http=host:port https=host:port socks=host:port")
	flag.StringVar(&off, "off", "", "Turn off services: http https socks all")
	flag.BoolVar(&setIgnore, "set-ignore", false, "Set default ignore")
	flag.BoolVar(&listenQuit, "listen-quit", false, "Listen quit do off all")
	flag.Parse()
}

func main() {
	switch {
	case on != "":
		onFunc()
	case off != "":
		offFunc()
	case setIgnore:
		setIgnoreFunc()
	case listenQuit:
		listenQuitFunc()
	default:
		flag.Usage()
		os.Exit(1)
		return
	}
}

func onFunc() {
	params := strings.Split(on, " ")
	if len(params) == 0 {
		fmt.Printf("Error: fail to on, no parameters.\n")
		return
	}

	var services []Service
	for _, s := range params {
		matches := serviceRegex.FindStringSubmatch(s)
		if matches == nil || len(matches) < 4 {
			fmt.Printf("Warn: fail to on, invalid parameter(%s).\n", s)
			continue
		}
		protocol := strings.ToUpper(matches[1])
		if protocol != "HTTP" && protocol != "HTTPS" && protocol != "SOCKS" {
			fmt.Printf("Warn: fail to on, invalid protocol(%s) [http https socks].\n", s)
			continue
		}
		host := matches[2]
		if net.ParseIP(host) == nil {
			fmt.Printf("Warn: fail to on, invalid ip(%s).\n", s)
			continue
		}
		port, err := strconv.Atoi(matches[3])
		if err != nil || !(port > 0 && port <= 65535) {
			fmt.Printf("Warn: fail to on, invalid port(%s).\n", s)
			continue
		}
		services = append(services, Service{
			Original: s,
			Protocol: protocol,
			Addr: sysproxy.Addr{
				Host: host,
				Port: port,
			},
		})
	}
	if len(services) == 0 {
		fmt.Printf("Error: fail to on, no available parameters.\n")
		return
	}

	for _, service := range services {
		var err error
		switch service.Protocol {
		case "HTTP":
			err = sysproxy.OnHttp(service.Addr)
		case "HTTPS":
			err = sysproxy.OnHttps(service.Addr)
		case "SOCKS":
			err = sysproxy.OnSocks(service.Addr)
		}
		if err != nil {
			fmt.Printf("Warn: fail to on, reason(%s): %s.\n", service.Original, err.Error())
		} else {
			fmt.Printf("Info: success to on, %s.\n", service.Original)
		}
	}
}

func offFunc() {
	params := strings.Split(off, " ")
	if len(params) == 0 {
		fmt.Printf("Error: fail to off, no parameters.\n")
		return
	}

	var protocols []string
	var containAll bool
	for _, s := range params {
		protocol := strings.ToUpper(s)
		if protocol != "HTTP" && protocol != "HTTPS" && protocol != "SOCKS" && protocol != "ALL" {
			fmt.Printf("Warn: fail to off, invalid protocol(%s) [http https socks].\n", s)
			continue
		}
		if protocol == "ALL" {
			containAll = true
		}
		protocols = append(protocols, protocol)
	}

	if containAll {
		if err := sysproxy.OffAll(); err != nil {
			fmt.Printf("Warn: fail to off, reason(all): %s.\n", err.Error())
		} else {
			fmt.Printf("Info: success to off, all.\n")
		}
		return
	}
	for _, protocol := range protocols {
		var err error
		switch protocol {
		case "HTTP":
			err = sysproxy.OffHttp()
		case "HTTPS":
			err = sysproxy.OffHttps()
		case "SOCKS":
			err = sysproxy.OffSocks()
		}
		if err != nil {
			fmt.Printf("Warn: fail to off, reason(%s): %s.\n", protocol, err.Error())
		} else {
			fmt.Printf("Info: success to off %s.\n", protocol)
		}
	}
}

func setIgnoreFunc() {
	if err := sysproxy.SetIgnore(sysproxy.DefaultIgnores); err != nil {
		fmt.Printf("Warn: fail to set ignore, reason: %s.\n", err.Error())
	} else {
		fmt.Printf("Info: success to set ignore.\n")
	}
}

func listenQuitFunc() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-c
	if err := sysproxy.OffAll(); err != nil {
		fmt.Println("listen do off all failed:", err)
	} else {
		fmt.Println("listen do off all success")
	}
}
