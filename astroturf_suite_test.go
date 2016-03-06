package astroturf_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAstroturf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Astroturf Suite")
}
