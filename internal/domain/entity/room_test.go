package entity

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
)

func TestRoom_ShouldCreateARoomWhenDataAreValid(t *testing.T) {
	id := valueobject.NewId()
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	createdAt := valueobject.NewTimestamp()
	updateAt := valueobject.NewTimestamp()

	room, err := NewRoom(adminId, name, category)
	assert.NotNil(t, room)
	assert.Nil(t, err)
	assert.NotNil(t, room.Id())
	assert.NotNil(t, room.AdminId())
	assert.NotNil(t, room.Name())
	assert.NotNil(t, room.Category())
	assert.NotNil(t, room.CreatedAt())
	assert.NotNil(t, room.UpdatedAt())
	assert.Equal(t, adminId.Value(), room.AdminId().Value())
	assert.Equal(t, name.Value(), room.Name().Value())
	assert.Equal(t, category.Value(), room.Category().Value())

	room, err = NewRoomWith(id, adminId, name, category, createdAt, updateAt)
	assert.NotNil(t, room)
	assert.Nil(t, err)
	assert.NotNil(t, room.Id())
	assert.NotNil(t, room.AdminId())
	assert.NotNil(t, room.Name())
	assert.NotNil(t, room.Category())
	assert.NotNil(t, room.CreatedAt())
	assert.NotNil(t, room.UpdatedAt())
	assert.Equal(t, id.Value(), room.Id().Value())
	assert.Equal(t, adminId.Value(), room.AdminId().Value())
	assert.Equal(t, name.Value(), room.Name().Value())
	assert.Equal(t, category.Value(), room.Category().Value())
	assert.Equal(t, createdAt.Value(), room.CreatedAt().Value())
	assert.Equal(t, updateAt.Value(), room.UpdatedAt().Value())
}

func TestRoom_ShouldReturnAnErrorWhenCreateARoomWithInvalidData(t *testing.T) {
	id := valueobject.NewId()
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	createdAt := valueobject.NewTimestamp()
	updatedAt := valueobject.NewTimestamp()

	testCases := []struct {
		testName  string
		errNew    error
		errWith   error
		id        *valueobject.Id
		adminId   *valueobject.UserId
		name      *valueobject.RoomName
		category  *valueobject.RoomCategory
		createdAt *valueobject.Timestamp
		updatedAt *valueobject.Timestamp
	}{
		{
			"nil id",
			nil,
			validation.ErrRequiredRoomId,
			nil,
			adminId,
			name,
			category,
			createdAt,
			updatedAt,
		},
		{
			"nil admin id",
			validation.ErrRequiredRoomAdminId,
			validation.ErrRequiredRoomAdminId,
			id,
			nil,
			name,
			category,
			createdAt,
			updatedAt,
		},
		{
			"nil room name",
			validation.ErrRequiredRoomName,
			validation.ErrRequiredRoomName,
			id,
			adminId,
			nil,
			category,
			createdAt,
			updatedAt,
		},
		{
			"nil room category",
			validation.ErrRequiredRoomCategory,
			validation.ErrRequiredRoomCategory,
			id,
			adminId,
			name,
			nil,
			createdAt,
			updatedAt,
		},
		{
			"nil created at",
			nil,
			validation.ErrRequiredRoomCreatedAt,
			id,
			adminId,
			name,
			category,
			nil,
			updatedAt,
		},
		{
			"nil updated at",
			nil,
			validation.ErrRequiredRoomUpdatedAt,
			id,
			adminId,
			name,
			category,
			createdAt,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			_, err := NewRoom(
				test.adminId,
				test.name,
				test.category,
			)
			assert.Equal(t, test.errNew, err)

			_, err = NewRoomWith(
				test.id,
				test.adminId,
				test.name,
				test.category,
				test.createdAt,
				test.updatedAt,
			)
			assert.Equal(t, test.errWith, err)
		})
	}
}

func TestRoom_ShouldUpdateRoomWhenDateAreValid(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := NewRoom(adminId, name, category)

	updatedAt := room.UpdatedAt()
	newName, _ := valueobject.NewRoomNameWith("Rust")
	err := room.UpdateName(newName)
	assert.Nil(t, err)
	assert.Equal(t, newName.Value(), room.Name().Value())
	assert.True(t, room.updatedAt.Value().After(updatedAt.Value()))

	updatedAt = room.UpdatedAt()
	newCategory, _ := valueobject.NewRoomCategoryWith("Science")
	err = room.UpdateCategory(newCategory)
	assert.Nil(t, err)
	assert.Equal(t, newCategory.Value(), room.Category().Value())
	assert.True(t, room.updatedAt.Value().After(updatedAt.Value()))
}

func TestRoom_ShouldReturnAnErrorWhenUpdateRoomWithInvalidData(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := NewRoom(adminId, name, category)

	updatedAt := room.UpdatedAt()
	err := room.UpdateName(nil)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrRequiredRoomName)
	assert.Equal(t, name.Value(), room.Name().Value())
	assert.True(t, room.updatedAt.Value().Equal(updatedAt.Value()))

	updatedAt = room.UpdatedAt()
	err = room.UpdateCategory(nil)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrRequiredRoomCategory)
	assert.Equal(t, category.Value(), room.Category().Value())
	assert.True(t, room.updatedAt.Value().Equal(updatedAt.Value()))
}
