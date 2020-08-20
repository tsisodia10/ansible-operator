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

![Screenshot](/a.png)

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

## Running AnsibleRunner Operator

### Pre Requisites 
1. Install docker 
```
$ sudo dnf config-manager --add-repo=https://download.docker.com/linux/centos/docker-ce.repo
```
```
$ sudo dnf install docker-ce-3:18.09.1-3.el7
```
2. Running OpenShift cluster

### Installation
1. You can pull the `AnsibleRunner Operator` image from -
```
docker pull quay.io/tsisodia/ansiblerunner-operator
```
2. Clone git repository
```
https://github.com/tsisodia10/ansible-operator
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

Configuration data can be consumed in pods in a variety of ways. A ConfigMap can be used to:

1. Populate the value of environment variables.
2. Set command-line arguments in a container.
3. Populate configuration files in a volume.
4. Both users and system components may store configuration data in a ConfigMap.

Ensure you create configMap and Secrets, for example - `password-secret.yaml`, `sshkey-secret.yaml`
```
oc create configmap <configmap_name> [options]
```

Example of ConfigMap with two environment variable:
```apiVersion: v1
kind: ConfigMap
metadata:
  name: special-config 
  namespace: default
data:
  special.how: very 
  special.type: charm 
```

Ensure you are currently scoped to the operator Namespace:
```
oc project operator
```

Deploy your AnsiblePlaybookRun Custom Resource to the live OpenShift Cluster:
```
oc create -f deploy/crds/ansible.konveyor.io_v1alpha1_ansibleplaybookrun_cr.yaml
```

### Running the Operator locally
Now we can test our logic by running our Operator outside the cluster. You can continue interacting with the OpenShift cluster by opening a new terminal window.
```
operator-sdk run local 
```


Verify the Ansible operator has created 1 job running Ansible Playbook :
```
oc get job
```
Verify in the OpenShift dashboard if the resources are created 
