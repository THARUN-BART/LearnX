package execution

import "context"

// Client is the interface for isolated code execution backends (Judge0, container sandboxes, etc.).
type Client interface {
	Execute(ctx context.Context, req ExecutionRequest) (*ExecutionResponse, error)
	ProviderName() string
}
