package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/codegangsta/cli"
	"github.com/miekg/dns"
)

func dnsCheck(name, ip, server string) (error, time.Duration) {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(name+".", dns.TypeA)
	r, t, err := c.Exchange(&m, server+":53")

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

func monitor(name, ip string, servers []string, interval time.Duration) {
	log.Printf("Monitoring...")

	namespace := "dd-dns-monitor"
	dd, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}
	dd.Namespace = namespace + "."

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			for _, server := range servers {
				tag := []string{server}
				err, t := dnsCheck(name, ip, server)
				if err != nil {
					log.Printf(err.Error()) //TODO change to debug
					dd.Count("error", 1, tag, 1)
				} else {
					dd.TimeInMilliseconds("time", float64(t/time.Millisecond), tag, 1)
				}
			}
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "dd-dns-monitor"
	app.Usage = "log DNS server failures to DataDog"
	app.Author = "Chris Kite"

	var name, ip, servers string
	var interval time.Duration

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
			Name:        "interval, t",
			Usage:       "interval in seconds to check at",
			Value:       15,
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
		monitor(name, ip, servers, interval*time.Second)
	}

	app.Run(os.Args)
}