package astroturf

import (
	"errors"
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/pivotal-golang/lager"
)

type backend struct {
	capacity      garden.Capacity
	containers    map[string]garden.Container
	graceTime     time.Duration
	containerLock sync.RWMutex
	logger        lager.Logger
}

func NewBackend(logger lager.Logger, memoryInBytes, diskInBytes, maxContainers uint64, graceTime time.Duration) *backend {
	containers := make(map[string]garden.Container)
	return &backend{
		capacity: garden.Capacity{
			MemoryInBytes: memoryInBytes,
			DiskInBytes:   diskInBytes,
			MaxContainers: maxContainers,
		},
		containers:    containers,
		graceTime:     graceTime,
		containerLock: sync.RWMutex{},
		logger:        logger,
	}
}

func (c *backend) GraceTime(container garden.Container) time.Duration {
	return c.graceTime
}

func (c *backend) Start() error { return nil }
func (c *backend) Stop()        {}
func (c *backend) Ping() error  { return nil }

func (c *backend) Capacity() (garden.Capacity, error) {
	return c.capacity, nil
}

func (c *backend) Create(spec garden.ContainerSpec) (garden.Container, error) {
	c.containerLock.Lock()
	defer c.containerLock.Unlock()

	_, ok := c.containers[spec.Handle]
	if ok {
		return nil, errors.New("handle already taken")
	}

	container := NewContainer(c.logger, spec)
	c.containers[spec.Handle] = container
	return container, nil
}

func (c *backend) Destroy(handle string) error {
	c.containerLock.Lock()
	defer c.containerLock.Unlock()

	if _, ok := c.containers[handle]; !ok {
		return errors.New("container does not exist")
	}

	delete(c.containers, handle)
	return nil
}

func (c *backend) Containers(properties garden.Properties) ([]garden.Container, error) {
	c.containerLock.RLock()
	defer c.containerLock.RUnlock()

	matchingContainers := []garden.Container{}
	for _, container := range c.containers {
		matched := true
		for k, v := range properties {
			value, err := container.Property(k)
			if err != nil {
				matched = false
				break
			}

			if v != value {
				matched = false
				break
			}
		}

		if matched {
			matchingContainers = append(matchingContainers, container)
		}
	}

	return matchingContainers, nil
}

func (c *backend) BulkInfo(handles []string) (map[string]garden.ContainerInfoEntry, error) {
	c.containerLock.RLock()
	defer c.containerLock.RUnlock()

	infos := make(map[string]garden.ContainerInfoEntry)
	for _, handle := range handles {
		container, ok := c.containers[handle]
		if !ok {
			continue
		}

		info, _ := container.Info()
		infos[handle] = garden.ContainerInfoEntry{info, nil}
	}

	return infos, nil
}

func (c *backend) BulkMetrics(handles []string) (map[string]garden.ContainerMetricsEntry, error) {
	c.containerLock.RLock()
	defer c.containerLock.RUnlock()

	metrics := make(map[string]garden.ContainerMetricsEntry)
	for _, handle := range handles {
		container, ok := c.containers[handle]
		if !ok {
			continue
		}

		metric, _ := container.Metrics()
		metrics[handle] = garden.ContainerMetricsEntry{metric, nil}
	}

	return metrics, nil
}

func (c *backend) Lookup(handle string) (garden.Container, error) {
	c.containerLock.RLock()
	defer c.containerLock.RUnlock()

	container, ok := c.containers[handle]
	if !ok {
		return nil, errors.New("container does not exist")
	}

	return container, nil
}
