package ansibleplaybookrun_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAnsibleplaybookrun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ansibleplaybookrun Suite")
}
