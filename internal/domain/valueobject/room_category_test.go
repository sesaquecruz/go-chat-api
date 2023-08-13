package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"

	"github.com/stretchr/testify/assert"
)

func TestRoomCategory_ShouldCreateANewCategoryWhenValueIsValid(t *testing.T) {
	categories := []string{
		"General", "Tech", "Game", "Book", "Movie", "Music", "Language", "Science",
	}

	for _, value := range categories {
		category, err := NewRoomCategoryWith(value)
		assert.NotNil(t, category)
		assert.Nil(t, err)
		assert.Equal(t, value, category.Value())
	}
}

func TestRoomCategory_ShouldReturnARequireRoomCategoryErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	category, err := NewRoomCategoryWith(value)
	assert.Nil(t, category)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomCategory)
}

func TestRoomCategory_ShouldReturnAInvalidRoomCategoryErrorWhenValueIsInvalid(t *testing.T) {
	value := "Other"
	category, err := NewRoomCategoryWith(value)
	assert.Nil(t, category)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrInvalidRoomCategory)
}
