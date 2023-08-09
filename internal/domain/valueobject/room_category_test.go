package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewCategoryWhenValueIsValid(t *testing.T) {
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

func TestShouldReturnARequireRoomCategoryErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	category, err := NewRoomCategoryWith(value)
	assert.Nil(t, category)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomCategory)
}

func TestShouldReturnAInvalidRoomCategoryErrorWhenValueIsInvalid(t *testing.T) {
	value := "Other"
	category, err := NewRoomCategoryWith(value)
	assert.Nil(t, category)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, ErrInvalidRoomCategory)
}
