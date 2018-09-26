package blueprint

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBlueprint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Blueprint Suite")
}
