package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/search"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type RoomGatewayTestSuite struct {
	suite.Suite
	ctx       context.Context
	container *test.PostgresContainer
	gateway   *RoomGateway
}

func (s *RoomGatewayTestSuite) SetupSuite() {
	ctx := context.Background()
	container, err := test.NewPostgresContainer(ctx, "file://../../../migrations")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", container.DSN)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	container.ClearDB()

	s.ctx = ctx
	s.container = container
	s.gateway = NewRoomGateway(db)
}

func (s *RoomGatewayTestSuite) TearDownSuite() {
	if err := s.container.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func TestRoomGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(RoomGatewayTestSuite))
}

func (s *RoomGatewayTestSuite) TestSave_ShouldSaveAndFindARoom() {
	defer s.container.ClearDB()
	t := s.T()

	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Resident Evil")
	category, _ := valueobject.NewRoomCategoryWith("Game")

	room, err := entity.NewRoom(adminId, name, category)
	assert.Nil(t, err)

	err = s.gateway.Save(s.ctx, room)
	assert.Nil(t, err)

	result, err := s.gateway.FindById(s.ctx, room.Id())
	assert.NotNil(t, result)
	assert.Nil(t, err)

	assert.Equal(t, room.Id().Value(), result.Id().Value())
	assert.Equal(t, room.AdminId().Value(), result.AdminId().Value())
	assert.Equal(t, room.Name().Value(), result.Name().Value())
	assert.Equal(t, room.Category().Value(), result.Category().Value())
	assert.Equal(t, room.CreatedAt().Value(), result.CreatedAt().Value())
	assert.Equal(t, room.UpdatedAt().Value(), result.CreatedAt().Value())
}

func (s *RoomGatewayTestSuite) TestFind_ShouldReturnAnErrorWhenFindANonexistentRoom() {
	defer s.container.ClearDB()
	t := s.T()

	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Resident Evil")
	category, _ := valueobject.NewRoomCategoryWith("Game")

	room, err := entity.NewRoom(adminId, name, category)
	assert.Nil(t, err)

	result, err := s.gateway.FindById(s.ctx, room.Id())
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func (s *RoomGatewayTestSuite) TestSearch_ShouldReturnRoomPages() {
	defer s.container.ClearDB()
	t := s.T()

	for i := 0; i < 10; i++ {
		adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
		name, _ := valueobject.NewRoomNameWith(fmt.Sprintf("Room %d", i))
		category, _ := valueobject.NewRoomCategoryWith("Game")
		room, _ := entity.NewRoom(adminId, name, category)
		s.gateway.Save(s.ctx, room)
	}

	for i := 0; i < 5; i++ {
		query, err := search.NewQuery(strconv.Itoa(i), "2", "", "")
		assert.Nil(t, err)

		page, err := s.gateway.Search(s.ctx, query)
		assert.Nil(t, err)
		assert.Equal(t, i, page.Page)
		assert.Equal(t, 2, page.Size)
		assert.Equal(t, int64(10), page.Total)
		assert.Equal(t, 2, len(page.Items))

		for j := i * 2; j < 2; j++ {
			assert.Equal(t, fmt.Sprintf("Room %d", j), page.Items[j].Name().Value())
		}
	}

	for i := 0; i < 5; i++ {
		query, err := search.NewQuery(strconv.Itoa(i), "2", "asc", "")
		assert.Nil(t, err)

		page, err := s.gateway.Search(s.ctx, query)
		assert.Nil(t, err)
		assert.Equal(t, i, page.Page)
		assert.Equal(t, 2, page.Size)
		assert.Equal(t, int64(10), page.Total)
		assert.Equal(t, 2, len(page.Items))

		for j := i * 2; j < 2; j++ {
			assert.Equal(t, fmt.Sprintf("Room %d", j), page.Items[j].Name().Value())
		}
	}

	for i := 4; i >= 0; i-- {
		query, err := search.NewQuery(strconv.Itoa(i), "2", "desc", "")
		assert.Nil(t, err)

		page, err := s.gateway.Search(s.ctx, query)
		assert.Nil(t, err)
		assert.Equal(t, i, page.Page)
		assert.Equal(t, 2, page.Size)
		assert.Equal(t, int64(10), page.Total)
		assert.Equal(t, 2, len(page.Items))

		for j := 0; j < 2; j++ {
			assert.Equal(t, fmt.Sprintf("Room %d", 9-i*2-j), page.Items[j].Name().Value())
		}
	}
}

func (s *RoomGatewayTestSuite) TestSearch_ShouldReturnRoomBySearchTerm() {
	defer s.container.ClearDB()
	t := s.T()

	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")

	name, _ := valueobject.NewRoomNameWith("Need for Speed Most Wanted")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, _ := entity.NewRoom(adminId, name, category)
	s.gateway.Save(s.ctx, room)

	name, _ = valueobject.NewRoomNameWith("Need for Speed Underground")
	category, _ = valueobject.NewRoomCategoryWith("Game")
	room, _ = entity.NewRoom(adminId, name, category)
	s.gateway.Save(s.ctx, room)

	name, _ = valueobject.NewRoomNameWith("Rust")
	category, _ = valueobject.NewRoomCategoryWith("Tech")
	room, _ = entity.NewRoom(adminId, name, category)
	s.gateway.Save(s.ctx, room)

	name, _ = valueobject.NewRoomNameWith("Go")
	category, _ = valueobject.NewRoomCategoryWith("Tech")
	room, _ = entity.NewRoom(adminId, name, category)
	s.gateway.Save(s.ctx, room)

	name, _ = valueobject.NewRoomNameWith("Java")
	category, _ = valueobject.NewRoomCategoryWith("Tech")
	room, _ = entity.NewRoom(adminId, name, category)
	s.gateway.Save(s.ctx, room)

	query, _ := search.NewQuery("", "", "desc", "for Speed")
	page, _ := s.gateway.Search(s.ctx, query)
	assert.Equal(t, 0, page.Page)
	assert.Equal(t, 10, page.Size)
	assert.Equal(t, int64(2), page.Total)
	assert.Equal(t, 2, len(page.Items))
	assert.Equal(t, "Need for Speed Underground", page.Items[0].Name().Value())
	assert.Equal(t, "Need for Speed Most Wanted", page.Items[1].Name().Value())

	query, _ = search.NewQuery("", "", "", "Tech")
	page, _ = s.gateway.Search(s.ctx, query)
	assert.Equal(t, 0, page.Page)
	assert.Equal(t, 10, page.Size)
	assert.Equal(t, int64(3), page.Total)
	assert.Equal(t, 3, len(page.Items))
	assert.Equal(t, "Go", page.Items[0].Name().Value())
	assert.Equal(t, "Java", page.Items[1].Name().Value())
	assert.Equal(t, "Rust", page.Items[2].Name().Value())
}

func (s *RoomGatewayTestSuite) TestUpdate_ShouldUpdateARoom() {
	defer s.container.ClearDB()
	t := s.T()

	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Resident Evil")
	category, _ := valueobject.NewRoomCategoryWith("Game")

	room, err := entity.NewRoom(adminId, name, category)
	assert.Nil(t, err)

	err = s.gateway.Save(s.ctx, room)
	assert.Nil(t, err)

	newName, _ := valueobject.NewRoomNameWith("Rust")
	newCategory, _ := valueobject.NewRoomCategoryWith("Tech")

	room.UpdateName(newName)
	room.UpdateCategory(newCategory)

	err = s.gateway.Update(s.ctx, room)
	assert.Nil(t, err)

	result, err := s.gateway.FindById(s.ctx, room.Id())
	assert.NotNil(t, result)
	assert.Nil(t, err)

	assert.Equal(t, room.Id().Value(), result.Id().Value())
	assert.Equal(t, room.AdminId().Value(), result.AdminId().Value())
	assert.Equal(t, room.Name().Value(), result.Name().Value())
	assert.Equal(t, room.Category().Value(), result.Category().Value())
	assert.Equal(t, room.CreatedAt().Value(), result.CreatedAt().Value())
	assert.Equal(t, room.UpdatedAt().Value(), result.UpdatedAt().Value())
}

func (s *RoomGatewayTestSuite) TestDelete_ShouldDeleteARoom() {
	defer s.container.ClearDB()
	t := s.T()

	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Resident Evil")
	category, _ := valueobject.NewRoomCategoryWith("Game")

	room, err := entity.NewRoom(adminId, name, category)
	assert.Nil(t, err)

	err = s.gateway.Save(s.ctx, room)
	assert.Nil(t, err)

	err = s.gateway.Delete(s.ctx, room.Id())
	assert.Nil(t, err)

	result, err := s.gateway.FindById(s.ctx, room.Id())
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
