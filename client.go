package astroturf

import (
	"errors"

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
	if _, ok := c.containers[handle]; !ok {
		return errors.New("container does not exist")
	}

	delete(c.containers, handle)
	return nil
}

func (c *client) Containers(properties garden.Properties) ([]garden.Container, error) {
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

func (c *client) BulkInfo(handles []string) (map[string]garden.ContainerInfoEntry, error) {
	infos := make(map[string]garden.ContainerInfoEntry)
	for _, handle := range handles {
		container, ok := c.containers[handle]
		if !ok {
			continue
		}

		info, err := container.Info()
		infos[handle] = garden.ContainerInfoEntry{info, garden.NewError(err.Error())}
	}

	return infos, nil
}

func (c *client) BulkMetrics(handles []string) (map[string]garden.ContainerMetricsEntry, error) {
	metrics := make(map[string]garden.ContainerMetricsEntry)
	for _, handle := range handles {
		container, ok := c.containers[handle]
		if !ok {
			continue
		}

		metric, err := container.Metrics()
		metrics[handle] = garden.ContainerMetricsEntry{metric, garden.NewError(err.Error())}
	}

	return metrics, nil
}

func (c *client) Lookup(handle string) (garden.Container, error) {
	container, ok := c.containers[handle]
	if !ok {
		return nil, errors.New("container does not exist")
	}

	return container, nil
}
