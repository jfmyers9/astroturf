package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/cloudfoundry-incubator/garden/server"
	"github.com/jfmyers9/astroturf"
	"github.com/pivotal-golang/lager"
)

var listenNetwork = flag.String(
	"listenNetwork",
	"unix",
	"how to listen on the address (unix, tcp, etc.)",
)

var listenAddr = flag.String(
	"listenAddress",
	"/tmp/garden.sock",
	"address to listen on",
)

var containerGraceTime = flag.Duration(
	"containerGraceTime",
	0,
	"time after which to destroy idle containers",
)

var memoryInBytes = flag.Uint64(
	"memoryInBytes",
	0,
	"Total memory capacity in bytes",
)

var diskInBytes = flag.Uint64(
	"diskInBytes",
	0,
	"Total disk capacity in bytes",
)

var maxContainers = flag.Uint64(
	"maxContainers",
	0,
	"Maximum number of containers that can be created",
)

func main() {
	flag.Parse()

	graceTime := *containerGraceTime

	logger, _ := cf_lager.New("astroturf")
	logger.Info("starting")

	backend := astroturf.NewBackend(*memoryInBytes, *diskInBytes, *maxContainers, *containerGraceTime)

	gardenServer := server.New(*listenNetwork, *listenAddr, graceTime, backend, logger)
	err := gardenServer.Start()
	if err != nil {
		logger.Fatal("failed-to-start-server", err)
	}

	signals := make(chan os.Signal, 1)
	go func() {
		<-signals
		gardenServer.Stop()
		os.Exit(0)
	}()

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	logger.Info("started", lager.Data{
		"network": *listenNetwork,
		"addr":    *listenAddr,
	})

	select {}
}
