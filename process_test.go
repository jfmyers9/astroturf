package astroturf_test

import (
	"time"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/jfmyers9/astroturf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/clock/fakeclock"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Process", func() {
	var (
		logger    *lagertest.TestLogger
		fakeClock *fakeclock.FakeClock
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("process")
		fakeClock = fakeclock.NewFakeClock(time.Now())
	})

	Context("Wait", func() {
		Context("when no valid process result is specified", func() {
			It("exits immediately and successfully", func() {
				spec := garden.ProcessSpec{}
				process, err := astroturf.NewProcess(logger, spec, fakeClock)
				Expect(err).NotTo(HaveOccurred())

				exitCode, err := process.Wait()
				Expect(err).NotTo(HaveOccurred())
				Expect(exitCode).To(Equal(0))
			})
		})

		Context("when a valid process result is specified", func() {
			It("waits the duration and returns the exit code specified", func() {
				processResult := `{"duration_in_seconds": 120, "exit_code": 27}`
				spec := garden.ProcessSpec{Path: processResult}

				process, err := astroturf.NewProcess(logger, spec, fakeClock)
				Expect(err).NotTo(HaveOccurred())

				waitChan := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					exitCode, err := process.Wait()
					Expect(err).NotTo(HaveOccurred())
					Expect(exitCode).To(Equal(27))
					close(waitChan)
				}()

				Consistently(waitChan).ShouldNot(BeClosed())
				fakeClock.IncrementBySeconds(119)
				Consistently(waitChan).ShouldNot(BeClosed())
				fakeClock.IncrementBySeconds(121)
				Eventually(waitChan).Should(BeClosed())
			})
		})
	})

	Describe("Signal", func() {
		It("causes a waiting process to exit", func() {

		})
	})

	Describe("SetTTY", func() {
		It("does nothing", func() {
			process, err := astroturf.NewProcess(logger, garden.ProcessSpec{}, fakeClock)
			Expect(err).NotTo(HaveOccurred())

			err = process.SetTTY(garden.TTYSpec{})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
