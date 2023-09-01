package database

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var postgresRoomRepository, _ = services.NewPostgresContainer(context.Background(), "file://../../../")

type RoomPostgresRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	repository *RoomPostgresRepository
}

func (s *RoomPostgresRepositoryTestSuite) SetupSuite() {
	postgresRoomRepository.Clear()

	db := PostgresConnection(&config.DatabaseConfig{
		Host:     postgresRoomRepository.Host,
		Port:     postgresRoomRepository.Port,
		User:     postgresRoomRepository.User,
		Password: postgresRoomRepository.Password,
		Name:     postgresRoomRepository.Name,
	})

	s.ctx = context.Background()
	s.repository = NewRoomPostgresRepository(db)
}

func (s *RoomPostgresRepositoryTestSuite) TearDownSuite() {
	if err := postgresRoomRepository.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func TestRoomPostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RoomPostgresRepositoryTestSuite))
}

func (s *RoomPostgresRepositoryTestSuite) TestShouldSaveAndFindARoom() {
	defer postgresRoomRepository.Clear()
	t := s.T()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room := entity.NewRoom(adminId, name, category)

	err := s.repository.Save(s.ctx, room)
	assert.Nil(t, err)

	result, err := s.repository.FindById(s.ctx, room.Id())
	assert.NotNil(t, result)
	assert.Nil(t, err)

	assert.Equal(t, room.Id().Value(), result.Id().Value())
	assert.Equal(t, room.AdminId().Value(), result.AdminId().Value())
	assert.Equal(t, room.Name().Value(), result.Name().Value())
	assert.Equal(t, room.Category().Value(), result.Category().Value())
	assert.Equal(t, room.CreatedAt().Value(), result.CreatedAt().Value())
	assert.Equal(t, room.UpdatedAt().Value(), result.CreatedAt().Value())
}

func (s *RoomPostgresRepositoryTestSuite) TestShouldReturnAnErrorWhenRoomDoesNotExist() {
	defer postgresRoomRepository.Clear()
	t := s.T()

	result, err := s.repository.FindById(s.ctx, valueobject.NewId())
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFoundRoom)
}

func (s *RoomPostgresRepositoryTestSuite) TestShouldReturnARoomPageInAscendingOrder() {
	defer postgresRoomRepository.Clear()
	t := s.T()

	total := 10
	pageSize := 2

	for i := 0; i < total; i++ {
		adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
		name, _ := valueobject.NewRoomNameWith(fmt.Sprintf("A Game %d", i))
		category, _ := valueobject.NewRoomCategoryWith("Game")
		room := entity.NewRoom(adminId, name, category)
		s.repository.Save(s.ctx, room)
	}

	for page := 0; page < 5; page++ {
		query, err := pagination.NewQuery(strconv.Itoa(page), strconv.Itoa(pageSize), "asc", "")
		assert.NotNil(t, query)
		assert.Nil(t, err)

		result, err := s.repository.Search(s.ctx, query)
		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, page, result.Page)
		assert.Equal(t, pageSize, result.Size)
		assert.Equal(t, int64(total), result.Total)
		assert.Equal(t, pageSize, len(result.Items))

		for pageItem := 0; pageItem < pageSize; pageItem++ {
			expectedName := fmt.Sprintf("A Game %d", page*pageSize+pageItem)
			item := result.Items[pageItem]
			assert.Equal(t, expectedName, item.Name().Value())
		}
	}
}

func (s *RoomPostgresRepositoryTestSuite) TestShouldReturnARoomPageInDescendingOrder() {
	defer postgresRoomRepository.Clear()
	t := s.T()

	total := 10
	pageSize := 2

	for i := 0; i < total; i++ {
		adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
		name, _ := valueobject.NewRoomNameWith(fmt.Sprintf("A Game %d", i))
		category, _ := valueobject.NewRoomCategoryWith("Game")
		room := entity.NewRoom(adminId, name, category)
		s.repository.Save(s.ctx, room)
	}

	for page := 0; page < 5; page++ {
		query, err := pagination.NewQuery(strconv.Itoa(page), strconv.Itoa(pageSize), "desc", "")
		assert.NotNil(t, query)
		assert.Nil(t, err)

		result, err := s.repository.Search(s.ctx, query)
		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, page, result.Page)
		assert.Equal(t, pageSize, result.Size)
		assert.Equal(t, int64(total), result.Total)
		assert.Equal(t, pageSize, len(result.Items))

		for pageItem := 0; pageItem < pageSize; pageItem++ {
			expectedName := fmt.Sprintf("A Game %d", (total-1)-(page*pageSize+pageItem))
			item := result.Items[pageItem]
			assert.Equal(t, expectedName, item.Name().Value())
		}
	}
}

func (s *RoomPostgresRepositoryTestSuite) TestShouldReturnARoomPageFilteredBySearchTerm() {
	defer postgresRoomRepository.Clear()
	t := s.T()

	rooms := []struct {
		adminId  string
		name     string
		category string
	}{
		{
			"auth0|64c8457bb160e37c8c34533b",
			"Need for Speed Most Wanted",
			"Game",
		},
		{
			"auth0|64c8457bb160e37c8c34533c",
			"Need for Speed Underground",
			"Game",
		},
		{
			"auth0|64c8457bb160e37c8c34533d",
			"Rust",
			"Tech",
		},
		{
			"auth0|64c8457bb160e37c8c34533e",
			"Go",
			"Tech",
		},
		{
			"auth0|64c8457bb160e37c8c34533f",
			"Java",
			"Tech",
		},
		{
			"auth0|64c8457bb160e37c8c34533g",
			"Python",
			"Tech",
		},
	}

	for i := 0; i < len(rooms); i++ {
		adminId, _ := valueobject.NewUserIdWith(rooms[i].adminId)
		name, _ := valueobject.NewRoomNameWith(rooms[i].name)
		category, _ := valueobject.NewRoomCategoryWith(rooms[i].category)
		room := entity.NewRoom(adminId, name, category)
		s.repository.Save(s.ctx, room)
	}

	query, _ := pagination.NewQuery("0", "10", "desc", "for Speed")
	page, _ := s.repository.Search(s.ctx, query)
	assert.Equal(t, 0, page.Page)
	assert.Equal(t, 10, page.Size)
	assert.Equal(t, int64(2), page.Total)
	assert.Equal(t, 2, len(page.Items))
	assert.Equal(t, "Need for Speed Underground", page.Items[0].Name().Value())
	assert.Equal(t, "Need for Speed Most Wanted", page.Items[1].Name().Value())

	query, _ = pagination.NewQuery("1", "2", "asc", "Tech")
	page, _ = s.repository.Search(s.ctx, query)
	assert.Equal(t, 1, page.Page)
	assert.Equal(t, 2, page.Size)
	assert.Equal(t, int64(4), page.Total)
	assert.Equal(t, 2, len(page.Items))
	assert.Equal(t, "Python", page.Items[0].Name().Value())
	assert.Equal(t, "Rust", page.Items[1].Name().Value())
}

func (s *RoomPostgresRepositoryTestSuite) TestShouldUpdateARoom() {
	defer postgresRoomRepository.Clear()
	t := s.T()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Resident Evil")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room := entity.NewRoom(adminId, name, category)

	err := s.repository.Save(s.ctx, room)
	assert.Nil(t, err)

	newName, _ := valueobject.NewRoomNameWith("Rust")
	newCategory, _ := valueobject.NewRoomCategoryWith("Tech")

	room.UpdateName(newName)
	room.UpdateCategory(newCategory)

	err = s.repository.Update(s.ctx, room)
	assert.Nil(t, err)

	result, err := s.repository.FindById(s.ctx, room.Id())
	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, room.Id().Value(), result.Id().Value())
	assert.Equal(t, room.AdminId().Value(), result.AdminId().Value())
	assert.Equal(t, room.Name().Value(), result.Name().Value())
	assert.Equal(t, room.Category().Value(), result.Category().Value())
	assert.Equal(t, room.CreatedAt().Value(), result.CreatedAt().Value())
	assert.Equal(t, room.UpdatedAt().Value(), result.UpdatedAt().Value())
}

func (s *RoomPostgresRepositoryTestSuite) TestShouldDeleteARoom() {
	defer postgresRoomRepository.Clear()
	t := s.T()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Resident Evil")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room := entity.NewRoom(adminId, name, category)

	err := s.repository.Save(s.ctx, room)
	assert.Nil(t, err)

	result, err := s.repository.FindById(s.ctx, room.Id())
	assert.NotNil(t, result)
	assert.Nil(t, err)

	err = s.repository.Delete(s.ctx, room.Id())
	assert.Nil(t, err)

	result, err = s.repository.FindById(s.ctx, room.Id())
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFoundRoom)
}
