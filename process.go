package astroturf

import (
	"github.com/cloudfoundry-incubator/garden"
	"github.com/nu7hatch/gouuid"
)

type process struct {
	id string
}

func NewProcess() (*process, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &process{
		id: id.String(),
	}, nil
}

func (p *process) ID() string {
	return p.id
}

func (p *process) Wait() (int, error) {
	return 0, nil
}

func (p *process) SetTTY(spec garden.TTYSpec) error {
	return nil
}

func (p *process) Signal(signal garden.Signal) error {
	return nil
}
