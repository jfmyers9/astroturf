package astroturf

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/pivotal-golang/lager"
)

type container struct {
	handle          string
	bandwidthLimits garden.BandwidthLimits
	cpuLimits       garden.CPULimits
	diskLimits      garden.DiskLimits
	memoryLimits    garden.MemoryLimits
	properties      garden.Properties
	processes       map[string]garden.Process
	processLock     sync.RWMutex
	logger          lager.Logger
}

func NewContainer(logger lager.Logger, handle string) *container {
	return &container{
		handle:      handle,
		properties:  garden.Properties{},
		processes:   map[string]garden.Process{},
		processLock: sync.RWMutex{},
		logger:      logger,
	}
}

func (c *container) Handle() string {
	return c.handle
}

func (c *container) Stop(kill bool) error { return nil }

func (c *container) Info() (garden.ContainerInfo, error) {
	return garden.ContainerInfo{
		ExternalIP: "localhost",
	}, nil
}

func (c *container) StreamIn(spec garden.StreamInSpec) error                    { return nil }
func (c *container) StreamOut(spec garden.StreamOutSpec) (io.ReadCloser, error) { return nil, nil }

func (c *container) LimitBandwidth(limits garden.BandwidthLimits) error {
	c.bandwidthLimits = limits
	return nil
}

func (c *container) CurrentBandwidthLimits() (garden.BandwidthLimits, error) {
	return c.bandwidthLimits, nil
}

func (c *container) LimitCPU(limits garden.CPULimits) error {
	c.cpuLimits = limits
	return nil
}

func (c *container) CurrentCPULimits() (garden.CPULimits, error) {
	return c.cpuLimits, nil
}

func (c *container) LimitDisk(limits garden.DiskLimits) error {
	c.diskLimits = limits
	return nil
}

func (c *container) CurrentDiskLimits() (garden.DiskLimits, error) {
	return c.diskLimits, nil
}

func (c *container) LimitMemory(limits garden.MemoryLimits) error {
	c.memoryLimits = limits
	return nil
}

func (c *container) CurrentMemoryLimits() (garden.MemoryLimits, error) {
	return c.memoryLimits, nil
}

func (c *container) NetIn(hostPort, containerPort uint32) (uint32, uint32, error) {
	return hostPort, containerPort, nil
}

func (c *container) NetOut(netOutRule garden.NetOutRule) error {
	return nil
}

func (c *container) Run(processSpec garden.ProcessSpec, processIO garden.ProcessIO) (garden.Process, error) {
	logger := c.logger.Session("run", lager.Data{"spec": processSpec})
	logger.Info("starting")
	defer logger.Info("completed")

	c.processLock.Lock()
	defer c.processLock.Unlock()

	p, err := NewProcess(c.logger, processSpec)
	c.processes[p.ID()] = p
	return p, err
}

func (c *container) Attach(processID string, io garden.ProcessIO) (garden.Process, error) {
	c.processLock.RLock()
	defer c.processLock.RUnlock()

	p, ok := c.processes[processID]
	if !ok {
		return nil, errors.New("no process")
	}

	return p, nil
}

func (c *container) Metrics() (garden.Metrics, error) {
	return garden.Metrics{}, nil
}

func (c *container) SetGraceTime(graceTime time.Duration) error {
	return nil
}

func (c *container) Properties() (garden.Properties, error) {
	return c.properties, nil
}

func (c *container) Property(name string) (string, error) {
	var err error

	property, ok := c.properties[name]
	if !ok {
		err = errors.New("property does not exist")
	}
	return property, err
}

func (c *container) SetProperty(name string, value string) error {
	c.properties[name] = value
	return nil
}

func (c *container) RemoveProperty(name string) error {
	delete(c.properties, name)
	return nil
}
