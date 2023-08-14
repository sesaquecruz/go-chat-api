package database

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
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

func (s *RoomGatewayTestSuite) TestShouldSaveAndFindARoom() {
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

func (s *RoomGatewayTestSuite) TestShouldReturnAnErrorWhenFindANonexistentRoom() {
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

func TestRoomGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(RoomGatewayTestSuite))
}
