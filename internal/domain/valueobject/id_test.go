package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestId_ShouldCreateAIdWhenValueIsValid(t *testing.T) {
	id := NewId()
	_, err := uuid.Parse(id.Value())
	assert.Nil(t, err)

	value := uuid.NewString()
	id, err = NewIdWith(value)
	assert.NotNil(t, id)
	assert.Nil(t, err)
	assert.Equal(t, value, id.Value())
}

func TestId_ShouldReturnAValidationErrorWhenValueIsInvalid(t *testing.T) {
	testCases := []struct {
		test  string
		value string
		err   error
	}{
		{
			"empty value",
			"",
			ErrRequiredId,
		},
		{
			"invalid value",
			"kj12389013kjfdsf9819023jkhfjds",
			ErrInvalidId,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			id, err := NewIdWith(tc.value)
			assert.Nil(t, id)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
			assert.IsType(t, validation.ValidationError(""), err)
		})
	}
}
