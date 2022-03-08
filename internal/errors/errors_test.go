package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExternal(t *testing.T) {
	tests := map[string]struct {
		err      error
		external bool
	}{
		"External error": {
			err:      &ExternalError{},
			external: true,
		},
		"Non-external error": {
			err:      nil,
			external: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			external := IsExternal(test.err)
			assert.Equal(external, test.external)
		})
	}
}
