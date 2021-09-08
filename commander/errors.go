package commander

import "fmt"

type CannotBootstrapError struct {
	reason string
}

func NewCannotBootstrapError(reason string) *CannotBootstrapError {
	return &CannotBootstrapError{reason}
}

func (c CannotBootstrapError) Error() string {
	return fmt.Sprintf("cannot bootstrap: %s", c.reason)
}

type InconsistentChainIDError struct {
	CannotBootstrapError
}

func NewInconsistentChainIDError(chainIDSource string) *InconsistentChainIDError {
	reason := fmt.Sprintf("chain ID conflict between config and %s", chainIDSource)
	return &InconsistentChainIDError{
		CannotBootstrapError: *NewCannotBootstrapError(reason),
	}
}
