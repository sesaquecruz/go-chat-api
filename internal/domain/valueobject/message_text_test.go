package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestMessageText_ShouldCreateAMessageTextWhenValueIsValid(t *testing.T) {
	text := "A simple message"
	messageText, err := NewMessageTextWith(text)
	assert.NotNil(t, messageText)
	assert.Nil(t, err)
	assert.Equal(t, text, messageText.Value())
}

func TestMessageText_ShouldReturnAValidationErrorWhenValueIsInvalid(t *testing.T) {
	testCases := []struct {
		test string
		text string
		err  error
	}{
		{
			"empty text",
			"",
			ErrRequiredMessageText,
		},
		{
			"blank text",
			"     ",
			ErrRequiredMessageText,
		},
		{
			"invalid text size",
			"dfafoiuereioqurdfjaiodfuweru19812321jioudfaudf 123u123u123ujkjfdsfu123u1239udfjsdlkfj12310  293213*12",
			ErrInvalidMessageText,
		},
	}

	assert.Equal(t, len(testCases[2].text), 101)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			messageText, err := NewMessageTextWith(tc.text)
			assert.Nil(t, messageText)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
			assert.IsType(t, validation.ValidationError(""), err)
		})
	}
}
