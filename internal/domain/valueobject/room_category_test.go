package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestRoomCategory_ShouldCreateARoomCategoryWhenValueIsValid(t *testing.T) {
	categories := []string{
		"General", "Tech", "Game", "Book", "Movie", "Music", "Language", "Science",
	}

	for _, value := range categories {
		roomCategory, err := NewRoomCategoryWith(value)
		assert.NotNil(t, roomCategory)
		assert.Nil(t, err)
		assert.Equal(t, value, roomCategory.Value())
	}
}

func TestRoomCategory_ShouldReturnAValidationErrorWhenValueIsInvalid(t *testing.T) {
	testCases := []struct {
		test  string
		value string
		err   error
	}{
		{
			"empty value",
			"",
			ErrRequiredRoomCategory,
		},
		{
			"invalid value",
			"Other",
			ErrInvalidRoomCategory,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			category, err := NewRoomCategoryWith(tc.value)
			assert.Nil(t, category)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
			assert.IsType(t, validation.ValidationError(""), err)
		})
	}
}
