package constants

// ResourceType describes the K8S types available
type ResourceType string

const (
	Deployments ResourceType = "deployments"
	Pods        ResourceType = "pods"
)
