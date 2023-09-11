package impl

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchRoomUseCase_ShouldReturnAPageWhenDataIsValid(t *testing.T) {
	ctx := context.Background()
	input := &usecase.SearchRoomUseCaseInput{
		Page:   "0",
		Size:   "2",
		Sort:   "asc",
		Search: "car",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.EXPECT().
		Search(mock.Anything, mock.Anything).
		Run(func(c context.Context, q *pagination.Query) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.Page, strconv.Itoa(q.Page()))
			assert.Equal(t, input.Size, strconv.Itoa(q.Size()))
			assert.Equal(t, strings.ToUpper(input.Sort), q.Sort())
			assert.Equal(t, strings.ToUpper(input.Search), q.Search())
		}).
		Return(pagination.NewPage[*entity.Room](0, 2, int64(10), []*entity.Room{}), nil).
		Once()

	useCase := NewSearchRoomUseCase(roomRepository)

	output, err := useCase.Execute(ctx, input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, 0, output.Page)
	assert.Equal(t, 2, output.Size)
	assert.Equal(t, int64(10), output.Total)
	assert.Equal(t, 0, len(output.Items))
}

func TestSearchRoomUseCase_ShouldReturnAnErrorWhenDataIsInvalid(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		test  string
		input *usecase.SearchRoomUseCaseInput
		err   error
	}{
		{
			"invalid page",
			&usecase.SearchRoomUseCaseInput{
				Page:   "-1",
				Size:   "2",
				Sort:   "asc",
				Search: "car",
			},
			pagination.ErrInvalidQueryPage,
		},
		{
			"invalid size",
			&usecase.SearchRoomUseCaseInput{
				Page:   "0",
				Size:   "0",
				Sort:   "asc",
				Search: "car",
			},
			pagination.ErrInvalidQuerySize,
		},
		{
			"invalid sort",
			&usecase.SearchRoomUseCaseInput{
				Page:   "0",
				Size:   "1",
				Sort:   "dfoierewr",
				Search: "car",
			},
			pagination.ErrInvalidQuerySort,
		},
		{
			"invalid search",
			&usecase.SearchRoomUseCaseInput{
				Page:   "0",
				Size:   "2",
				Sort:   "asc",
				Search: "dfaiuerewnvdiuoriewruuiwqeuqwe89123jladjsdasadiou23",
			},
			pagination.ErrInvalidQuerySearch,
		},
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)
	useCase := NewSearchRoomUseCase(roomRepository)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			output, err := useCase.Execute(ctx, tc.input)
			assert.Nil(t, output)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestSearchRoomUseCase_ShouldReturnAnErrorOnRepositoryError(t *testing.T) {
	ctx := context.Background()
	input := &usecase.SearchRoomUseCaseInput{
		Page:   "0",
		Size:   "2",
		Sort:   "asc",
		Search: "car",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.EXPECT().
		Search(mock.Anything, mock.Anything).
		Return(nil, errors.New("a repository error")).
		Once()

	useCase := NewSearchRoomUseCase(roomRepository)

	output, err := useCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "a repository error")
}
