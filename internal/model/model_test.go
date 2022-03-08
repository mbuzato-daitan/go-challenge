package model

import (
	"testing"

	"github.com/lucsky/cuid"
	"github.com/mtbuzato/go-challenge/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateID(t *testing.T) {
	tests := map[string]struct {
		id       string
		err      string
		external bool
	}{
		"Valid task": {
			id: cuid.New(),
		},
		"Invalid task ID": {
			id:       "",
			err:      "Invalid task ID.",
			external: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			err := ValidateID(test.id)
			if test.err != "" {
				assert.Equal(err.Error(), test.err)
				assert.True(errors.IsExternal(err))
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := map[string]struct {
		name     string
		err      string
		external bool
	}{
		"Valid task": {
			name: "Task 1",
		},
		"Invalid task name": {
			name:     "",
			err:      "Invalid task name.",
			external: true,
		},
		"Task name too long": {
			name:     "The name of this task is far too long to be accepted by the validator of the Task model so it should generate an error that describes what happened.",
			err:      "Task name is too long.",
			external: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			err := ValidateName(test.name)
			if test.err != "" {
				assert.Equal(err.Error(), test.err)
				assert.True(errors.IsExternal(err))
			} else {
				assert.NoError(err)
			}
		})
	}
}
