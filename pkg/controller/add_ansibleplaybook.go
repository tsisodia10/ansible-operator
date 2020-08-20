package controller

import (
	"github.com/ansible-operator/pkg/controller/ansibleplaybook"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ansibleplaybook.Add)
}
