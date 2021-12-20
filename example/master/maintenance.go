package main

import (
	"time"
)

func maintenanceInit() *time.Ticker {
	t := time.NewTicker(time.Duration(60) * time.Second)
	LogComponent("maintenance", "will run every 60 seconds")
	go performMaintenance(t)
	return t
}

func maintenanceShutdown(t *time.Ticker) {
	LogComponent("maintenance", "shutdown requested")
	t.Stop()
}

func performMaintenance(t *time.Ticker) {
	for range t.C {
		stale := cleanUpStaleServers()
		LogComponent("maintenance", "Cleaned up %d stale servers\n", stale)
	}
}

func cleanUpStaleServers() int {
	count := 0
	for k, v := range thisMaster.Service.Servers {
		if v.IsExpired(config.serverTimeout) {
			thisMaster.Lock()
			LogComponent("maintenance", "Removing server %s, last seen: %s", v.String(), v.LastSeen.Format(time.Stamp))
			delete(thisMaster.Service.Servers, k)
			delete(thisMaster.SolicitedServers, k)
			thisMaster.IPServiceCount[k]--
			if thisMaster.IPServiceCount[k] <= 0 {
				delete(thisMaster.IPServiceCount, k)
			}
			thisMaster.Unlock()
			count++
		}
	}
	return count
}
