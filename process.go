package astroturf

import (
	"encoding/json"
	"time"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
)

type process struct {
	id         string
	result     processResult
	exitStatus chan int
	signaled   chan garden.Signal
	logger     lager.Logger
	clock      clock.Clock
}

type processResult struct {
	Duration int `json:"duration_in_seconds"`
	ExitCode int `json:"exit_code"`
}

func NewProcess(logger lager.Logger, spec garden.ProcessSpec, clock clock.Clock) (*process, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	logger.Debug("garden-process-spec", lager.Data{"spec": spec})

	var result processResult
	err = json.Unmarshal([]byte(spec.Path), &result)
	if err != nil {
		logger.Error("failed-to-unmarshal-result", err)
		result = processResult{
			Duration: 0,
			ExitCode: 0,
		}
	}

	logger.Debug("created-result", lager.Data{"result": result})

	proc := &process{
		result:     result,
		id:         id.String(),
		exitStatus: make(chan int),
		signaled:   make(chan garden.Signal),
		logger:     logger,
		clock:      clock,
	}

	logger.Debug("created-process", lager.Data{"process": proc})

	go proc.run()

	return proc, nil
}

func (p *process) ID() string {
	return p.id
}

func (p *process) Wait() (int, error) {
	logger := p.logger.Session("waiting")
	logger.Info("starting")
	defer logger.Info("completed")

	return <-p.exitStatus, nil
}

func (p *process) SetTTY(spec garden.TTYSpec) error {
	return nil
}

func (p *process) Signal(signal garden.Signal) error {
	logger := p.logger.Session("signalling")
	logger.Info("starting")
	defer logger.Info("completed")

	p.signaled <- signal
	return nil
}

func (p *process) run() {
	logger := p.logger.Session("run", lager.Data{"result": p.result})
	logger.Info("starting")
	defer logger.Info("completed")

	timer := p.clock.NewTimer(time.Duration(p.result.Duration) * time.Second)

	select {
	case <-timer.C():
		p.exitStatus <- p.result.ExitCode
	case <-p.signaled:
		p.exitStatus <- 1
	}
}
