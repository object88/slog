package constants

import (
	"fmt"
)

// ResourceType describes the K8S types available
type ResourceType string

const (
	// Deployments represents the Kubernetes `deployment` resource
	Deployments ResourceType = "deployments"

	// Pods represents the Kubernetes `pods` resource
	Pods ResourceType = "pods"
)

// Validate ensures that the provided ResourceType is a valid value
func (rt ResourceType) Validate() error {
	switch rt {
	case Deployments, Pods:
		// Valid; nothing to do
		return nil
	default:
		return NewInvalidResourceTypeError(rt)
	}
}

// InvalidResourceTypeError is returned from ResourceType.Validate when the
// supplied resource type is not valid
type InvalidResourceTypeError struct {
	rt ResourceType
}

// NewInvalidResourceTypeError returns a new instance of the
// InvalidResourceTypeError struct
func NewInvalidResourceTypeError(resourceType ResourceType) *InvalidResourceTypeError {
	return &InvalidResourceTypeError{rt: resourceType}
}

func (irte *InvalidResourceTypeError) Error() string {
	return fmt.Sprintf("Resource type '%s' is not valid", string(irte.rt))
}
