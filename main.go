package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/miekg/dns"
)

func dnsCheck(name, ip, server string, timeout time.Duration) (error, time.Duration) {
	m := new(dns.Msg)
	m.SetQuestion(name+".", dns.TypeA)

	c := new(dns.Client)
	c.Dialer = &net.Dialer{
		Timeout: timeout,
	}
	r, t, err := c.Exchange(m, server+":53")

	if err != nil {
		return err, t
	}

	if len(r.Answer) == 0 {
		return err, t
	}

	for _, ans := range r.Answer {
		aRecord := ans.(*dns.A)
		addr := fmt.Sprintf("%s", aRecord.A)
		if addr != ip {
			err := fmt.Errorf("expected ip %s got %s", ip, addr)
			return err, t
		}
	}
	return nil, t
}

func monitor(name, ip string, servers []string, interval time.Duration, timeout time.Duration) {
	log.Printf("Monitoring...")

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			for _, server := range servers {
				err, t := dnsCheck(name, ip, server, timeout)
				ms := int32(t/time.Millisecond)
				if err != nil {
					log.Printf("DNS Check Error (%dms): %v", ms, err.Error())
				} else {
					log.Printf("DNS Check Success (%dms)", ms)
				}
			}
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "dd-dns-monitor"
	app.Usage = "log DNS server failures"
	app.Author = "Chris Kite"

	var name, ip, servers string
	var interval, timeout time.Duration

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "name, n",
			Usage:       "dns name to lookup",
			EnvVar:      "DNS_NAME",
			Destination: &name,
		},
		cli.StringFlag{
			Name:        "ip, i",
			Usage:       "ip address the dns name should resolve to",
			EnvVar:      "DNS_IP",
			Destination: &ip,
		},
		cli.StringFlag{
			Name:        "servers, s",
			Usage:       "comma-separated list of servers to monitor",
			EnvVar:      "DNS_SERVERS",
			Destination: &servers,
		},
		cli.DurationFlag{
			Name:        "timeout, t",
			Usage:       "dns query timeout",
			Value:       5 * time.Second,
			EnvVar:      "DNS_INTERVAL",
			Destination: &timeout,
		},
		cli.DurationFlag{
			Name:        "interval, l",
			Usage:       "interval in seconds to check at",
			Value:       500 * time.Millisecond,
			EnvVar:      "DNS_INTERVAL",
			Destination: &interval,
		},
	}

	app.Action = func(c *cli.Context) {
		if "" == name || "" == ip || "" == servers {
			cli.ShowAppHelp(c)
			return
		}

		servers := strings.Split(c.String("servers"), ",")
		monitor(name, ip, servers, interval, timeout)
	}

	app.Run(os.Args)
}
