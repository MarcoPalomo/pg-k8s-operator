package controller

import (
	"github.com/MarcoPalomo/pg-k8s-operator/pkg/controller/postgresuser"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, postgresuser.Add)
}
