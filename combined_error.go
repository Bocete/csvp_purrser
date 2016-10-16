package func_purrser

import (
	"strings"
)

type combinedError struct {
	e []error
}

func newCombinedError(e []error) combinedError {
	return combinedError{e: e}
}

func (ce combinedError) Error() string {
	descs := make([]string, len(ce.e))
	for i, e := range ce.e {
		descs[i] = e.Error()
	}
	return strings.Join(descs, "\n")
}
