package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/stretchr/testify/assert"
)

func TestMessageText_ShouldCreateAMessageText(t *testing.T) {
	text := "A simple message"
	messageText, err := NewMessageTextWith(text)
	assert.NotNil(t, messageText)
	assert.Nil(t, err)
}

func TestMessageText_ShouldReturnAnErrorWhenCreateAMessageTextWithInvalidDate(t *testing.T) {
	testCases := []struct {
		testName string
		text     string
		err      error
	}{
		{
			"empty case",
			"",
			validation.ErrRequiredMessageText,
		},
		{
			"blank case",
			"     ",
			validation.ErrRequiredMessageText,
		},
		{
			"max case",
			"dfafoiuereioqurdfjaiodfuweru19812321jioudfaudf 123u123u123ujkjfdsfu123u1239udfjsdlkfj12310  293213*12",
			validation.ErrSizeMessageText,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			messageText, err := NewMessageTextWith(test.text)
			assert.Nil(t, messageText)
			assert.ErrorIs(t, err, test.err)
		})
	}
}
