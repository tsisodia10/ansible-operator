// package ansibleplaybookrun_test

// import (
// 	"testing"

// 	ansiblev1alpha1 "github.com/ansible-operator/pkg/apis/ansible/v1alpha1"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )


// func TestAnsibleplaybookrun(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "Ansibleplaybookrun Suite")
// }

// var _ = Describe("Reconcile steps", func() {
// 	var (
// 		reconciler *ReconcileAnsiblePlaybookRun
// 		ap         *ansiblev1alpha1.AnsiblePlaybook
// 		apr        *ansiblev1alpha1.AnsiblePlaybookRun
// 	)

// 	Describe("Creating a Job", func() {
// 		BeforeEach(func() {
// 			ap.Spec.PlaybookName = "example-ansibleplaybook.yaml"
// 			ap.Spec.RepositoryType = "local"
// 			ap.Spec.RepositoryURL = "git"
// 			apr.Spec.ExtraVars = "extraVars.yaml"
// 			apr.Spec.Inventory = "inventory.yaml"
// 			apr.Spec.Password = "password.yaml"
// 			apr.Spec.SSHPrivateKey = "sshKey.yaml"
// 		})

// 		It("should fail to create if playbook is not provided: ", func() {
// 			ap.Spec.PlaybookName = nil

// 			result, err := reconciler.BuildJobSpec(ap, apr)

// 			Expect(result).To(BeNil())
// 			Expect(err).To(Not(BeNil()))
// 			Expect(err.Error()).To(Equal("Please provide playbook"))
// 		})

// 		It("should fail to create if repositoryType is not provided: ", func() {
// 			ap.Spec.RepositoryType = nil

// 			result, err := reconciler.BuildJobSpec(ap, apr)

// 			Expect(result).To(BeNil())
// 			Expect(err).To(Not(BeNil()))
// 			Expect(err.Error()).To(Equal("Please provide repository type"))
// 		})
// 	})
// })
