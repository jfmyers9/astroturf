package astroturf

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/garden"
)

type client struct {
	capacity   garden.Capacity
	containers map[string]garden.Container
}

func NewClient(memoryInBytes, diskInBytes, maxContainers uint64) *client {
	containers := make(map[string]garden.Container)
	return &client{
		capacity: garden.Capacity{
			MemoryInBytes: memoryInBytes,
			DiskInBytes:   diskInBytes,
			MaxContainers: maxContainers,
		},
		containers: containers,
	}
}

func (c *client) Ping() error { return nil }

func (c *client) Capacity() (garden.Capacity, error) {
	return c.capacity, nil
}

func (c *client) Create(spec garden.ContainerSpec) (garden.Container, error) {
	_, ok := c.containers[spec.Handle]
	if ok {
		return nil, errors.New("handle already taken")
	}

	container := NewContainer(spec.Handle)
	c.containers[spec.Handle] = container
	return container, nil
}

func (c *client) Destroy(handle string) error {
	delete(c.containers, handle)
	return nil
}

func (c *client) Containers(properties garden.Properties) ([]garden.Container, error) {
	return nil, nil
}

func (c *client) BulkInfo(handles []string) (map[string]garden.ContainerInfoEntry, error) {
	return nil, nil
}

func (c *client) BulkMetrics(handles []string) (map[string]garden.ContainerMetricsEntry, error) {
	return nil, nil
}

func (c *client) Lookup(handle string) (garden.Container, error) {
	return nil, nil
}

func (c *client) Start() error { return nil }

func (c *client) Stop() {}

func (c *client) GraceTime(container garden.Container) time.Duration {
	return 1 * time.Second
}
