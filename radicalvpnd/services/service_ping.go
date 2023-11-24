package service

import (
	"radicalvpnd/webapi"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func PingServers(server []webapi.Server) []webapi.Server {
	println("pinging", len(server), "servers")

	for i, s := range server {
		if s.Online {
			server[i] = pingServer(s)
		}
	}

	return server
}

func pingServer(s webapi.Server) webapi.Server {
	log.Debug("Pinging ", s.ExternaIp)

	pinger, err := probing.NewPinger(s.ExternaIp)

	if err != nil {
		log.Error("Failed to create pinger for", s.Hostname, err.Error())
		return s
	}

	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Timeout = time.Millisecond * 400
	pinger.Run()

	stats := pinger.Statistics()
	if stats.AvgRtt > 0 {
		latency := int(stats.AvgRtt / time.Millisecond)

		s.Latency = latency
	}

	return s
}
