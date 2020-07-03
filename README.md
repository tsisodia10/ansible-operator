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

