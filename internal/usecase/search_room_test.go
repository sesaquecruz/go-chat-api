package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/test/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSearchRoomUseCase_ShouldReturnAPageOfOutput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := SearchRoomUseCaseInput{
		Page:   "0",
		Size:   "2",
		Sort:   "asc",
		Search: "car",
	}

	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	repository.EXPECT().
		Search(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, q *pagination.Query) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.Page, strconv.Itoa(q.Page()))
			assert.Equal(t, input.Size, strconv.Itoa(q.Size()))
			assert.Equal(t, strings.ToUpper(input.Sort), q.Sort())
			assert.Equal(t, strings.ToUpper(input.Search), q.Search())
		}).
		Return(
			pagination.NewPage[*entity.Room](0, 2, int64(10), []*entity.Room{}),
			nil,
		).
		Times(1)

	useCase := NewSearchRoomUseCase(repository)
	output, err := useCase.Execute(ctx, &input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, 0, output.Page)
	assert.Equal(t, 2, output.Size)
	assert.Equal(t, int64(10), output.Total)
	assert.Equal(t, 0, len(output.Items))
}

func TestSearchRoomUseCase_ShouldReturnAnErrorWhenInputIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	useCase := NewSearchRoomUseCase(repository)

	testCases := []struct {
		input *SearchRoomUseCaseInput
		err   error
	}{
		{
			input: &SearchRoomUseCaseInput{
				Page:   "-1",
				Size:   "2",
				Sort:   "asc",
				Search: "car",
			},
			err: validation.ErrInvalidQueryPage,
		},
		{
			input: &SearchRoomUseCaseInput{
				Page:   "0",
				Size:   "0",
				Sort:   "asc",
				Search: "car",
			},
			err: validation.ErrInvalidQuerySize,
		},
		{
			input: &SearchRoomUseCaseInput{
				Page:   "0",
				Size:   "1",
				Sort:   "dfoierewr",
				Search: "car",
			},
			err: validation.ErrInvalidQuerySort,
		},
	}

	for _, test := range testCases {
		ouput, err := useCase.Execute(ctx, test.input)
		assert.Nil(t, ouput)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, test.err)
	}
}

func TestSearchRoomUseCase_ShouldReturnAnInternalErrorOnARepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := SearchRoomUseCaseInput{
		Page:   "0",
		Size:   "2",
		Sort:   "asc",
		Search: "car",
	}

	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	repository.EXPECT().
		Search(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("a repository error")).
		Times(1)

	useCase := NewSearchRoomUseCase(repository)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "a repository error")
}
