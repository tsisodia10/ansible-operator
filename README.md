# AnsibleRunner-Operator

## Use Cases & Design
The main use case is the ability to run an Ansible playbook against the virtual machine that is
migrated. For example, we want to remove the virtual machine from a load balancing pool. An
Ansible Playbook Run custom resource is created (cf. Custom Resources) and it contains all the
required information to run the playbook.

The following diagram [source] shows the corresponding flow.

The Ansible Controller creates three Persistent Volumes:
1. Project to install the playbook and its dependencies: roles, collections...
2. Inventory to store the inventory and variables, including credentials.
3. Artifacts to store the data produced by Ansible Runner: events, facts and output.

Then, the Ansible Controller creates the Ansible Runner job with the 3 PVs attached to it and
monitors the progress. The Ansible Runner pod copies the Ansible playbook repository, installs
the dependencies, creates the inventory and variables files, then calls the ansible-runner
command.

It might be worth creating a separate operator and depend on it to allow pre/post migration
hooks. The UI and admission controller could verify that the operator is deployed and enabled
when a Migration Plan is created.
For CAM, a hook-runner container image already exists and adds k8s and openshift libraries to
enable the respective modules. Because we will interact with VMware, RHV and maybe
OpenStack, we will have to extend this image to ship the required libraries.

## Custom Resources

### AnsiblePlaybook 
It represents an Ansible playbook specification to be run by the Ansible
Playbook controller.
- Repository
- Type
- URL
- Repository credentials secret
- Branch
- Playbook content
- Playbook file name

### AnsiblePlaybookRun 
It represents the execution of an Ansible Playbook.
AnsiblePlaybookRuns are how the AnsiblePlaybooks are executed ; they prepare the execution
and capture operational aspects of the AnsiblePlaybook execution such as events and progress,
so that the UI can represent them.
- AnsiblePlaybook (*AnsiblePlaybook)
- Inventory - List of hostnames / IP addresses extracted from provider
- Host(s) credentials ([]*Secrets)
- Extra vars - List of key / value pairs that are passed to the playbook at runtime

Ideally the Status of the CR will contain at least:

- State: pending, preparing, active, cleaning, finished
- Message: human readable status to display in the UI (Ansible task name)

## Building Operators

### Creating a new Project
Let's begin by creating a new project called operator : 

```
oc new-project operator
```

Let's now create a new directory in our `$GOPATH/src/` directory:

```
mkdir -p $GOPATH/src/github.com/redhat/
```

Navigate to the directory :

```
cd $GOPATH/src/github.com/redhat/
```

Create a new Go-based Operator SDK project for the AnsiblePlaybook:

```
operator-sdk new podset-operator --type=go
```

Navigate to the project root:

```
cd podset-operator
```

### Adding a new Custom API
Add two new Custom Resource Definition(CRD) APIs called `AnsiblePlaybook` and `AnsiblePlaybookRun`, with APIVersion ansible.konveyor.io/v1alpha1 and Kind `AnsiblePlaybook` and `AnsiblePlaybookRun` :
```
operator-sdk add api --api-version=ansible.konveyor.io/v1alpha1 --kind=AnsiblePLaybook
```

This will scaffold the `AnsiblePlaybook` and `AnsiblePlaybookRun` resource API under pkg/apis/app/v1alpha1/....

The Operator-SDK automatically creates the following manifests for you under the /deploy directory.

Custom Resource Definition
Custom Resource
Service Account
Role
RoleBinding
Deployment
Inspect the Custom Resource Definition manifest:
```
cat deploy/crds/app.example.com_ansibleplaybook_crd.yaml
```

### Defining the Spec and Status
Modify the AnsiblePlaybookRunSpec and AnsiblePlaybookRunStatus of the AnsiblePlaybook Custom Resource(CR) at 
```
go/src/github.com/redhat/ansible-operator/pkg/apis/app/v1alpha1/ansibleplaybookrun_types.go
```
Do the same woth AnsiblePlaybook API
It should look like the file below:
```
package v1alpha1

import (
	"context"

	kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	Pending   = "Pending"
	Preparing = "Preparing"
	Active    = "Active"
	Cleaning  = "Cleaning"
	Finished  = "Finished"
	Failed    = "Failed"
)

// AnsiblePlaybookRunSpec defines the desired state of AnsiblePlaybookRun
type AnsiblePlaybookRunSpec struct {
	AnsiblePlaybookRef *kapi.ObjectReference `json:"ansiblePlaybook,omitempty"`
	Inventory          string                `json:"inventory,omitempty"`
	HostCredential     string                `json:"hostCredential,omitempty"`
}

// AnsiblePlaybookRunStatus defines the observed state of AnsiblePlaybookRun
type AnsiblePlaybookRunStatus struct {

	Status string `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AnsiblePlaybookRun is the Schema for the ansibleplaybookruns API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ansibleplaybookruns,scope=Namespaced
type AnsiblePlaybookRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnsiblePlaybookRunSpec   `json:"spec,omitempty"`
	Status AnsiblePlaybookRunStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AnsiblePlaybookRunList contains a list of AnsiblePlaybookRun
type AnsiblePlaybookRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnsiblePlaybookRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnsiblePlaybookRun{}, &AnsiblePlaybookRunList{})
}

func (apr *AnsiblePlaybookRun) GetAnsiblePlaybook(client k8sclient.Client) (*AnsiblePlaybook, error) {

	//fmt.Printf("%+v\n", apr)
	if apr == nil {
		return nil, nil
	}
	object := AnsiblePlaybook{}
	// fmt.Printf("%+s\n", apr.Spec.AnsiblePlaybookRef.Namespace)
	// fmt.Printf("%+s\n", apr.Spec.AnsiblePlaybookRef.Name)
	err := client.Get(
		context.TODO(),
		types.NamespacedName{
			Namespace: apr.Spec.AnsiblePlaybookRef.Namespace,
			Name:      apr.Spec.AnsiblePlaybookRef.Name,
		},
		&object)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &object, err
}

```

After modifying the `*_types.go` file, always run the following command to update the generated code for that resource type:
```
operator-sdk generate k8s
```

We can also automatically update the CRD with OpenAPI v3 schema details based off the newly updated `*_types.go` file:
```
operator-sdk generate crds
```

Observe the CRD now reflects the spec.replicas and status.podNames OpenAPI v3 schema validation in the spec:
```
cat deploy/crds/ansible.konveyor.io_ansibleplaybook_crd.yaml
```

Deploy your `AnsiblePlaybook` Custom Resource Definition to the live OpenShift Cluster:
```
oc create -f deploy/crds/ansible.konveyor.io_ansibleplaybook_crd.yaml
```

Confirm the CRD was successfully created:
```
oc get crd ansibleplaybook.ansible.konveyor.io -o yaml
```

### Adding a new Controller
Add a new Controller to the project that will watch and reconcile the `AnsiblePlaybookRun` resource:
```
operator-sdk add controller --api-version=ansible.konveyor.io/v1alpha1 --kind=AnsiblePlaybookRun
```

This will scaffold a new Controller implementation under 
```
go/src/github.com/redhat/ansible-operator/pkg/controller/ansibleplaybookrun/ansibleplaybookrun_controller.go
```

### Customize the Operator Logic
Modify the PodSet controller logic at 
```
go/src/github.com/redhat/ansible-operator/pkg/controller/ansibleplaybookrun/ansibleplaybookrun_controller.go
```

### Running the Operator locally
Now we can test our logic by running our Operator outside the cluster. You can continue interacting with the OpenShift cluster by opening a new terminal window.
```
operator-sdk run local 
```

### Creating the Custom Resource 
In a new terminal, inspect the Custom Resource manifest:
```
cd $GOPATH/src/github.com/redhat/ansible-operator
cat deploy/crds/ansible.konveyor.io_v1alpha1_ansibleplaybookrun_cr.yaml
```
Ensure your kind: AnsiblePlaybookRun Custom Resource (CR) is updated with spec.replicas
```
apiVersion: ansible.konveyor.io/v1alpha1
kind: AnsiblePlaybookRun
metadata:
  name: example-ansibleplaybookrun
spec:
  # Add fields here
  ansiblePlaybook: 
    name: example-ansibleplaybook
    namespace: demo-test
  inventory: '{"all": {"hosts": "localhost"}}'
  hostCredential: '{"password": "b3BlcmF0b3I="}'
status: 
  status: pending
```

Ensure you are currently scoped to the operator Namespace:
```
oc project operator
```

Deploy your AnsiblePlaybookRun Custom Resource to the live OpenShift Cluster:
```
oc create -f deploy/crds/ansible.konveyor.io_v1alpha1_ansibleplaybookrun_cr.yaml
```

Verify the Ansible operator has created 1 job:
```
oc get job
```
