package entity

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
)

func TestRoom_ShouldCreateARoomWhenDateIsValid(t *testing.T) {
	id := valueobject.NewId()
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	createdAt := valueobject.NewTimestamp()
	updateAt := valueobject.NewTimestamp()
	var deleteAt *valueobject.Timestamp = nil

	room := NewRoom(adminId, name, category)
	assert.NotNil(t, room.Id())
	assert.Equal(t, adminId.Value(), room.AdminId().Value())
	assert.Equal(t, name.Value(), room.Name().Value())
	assert.Equal(t, category.Value(), room.Category().Value())
	assert.NotNil(t, room.CreatedAt())
	assert.NotNil(t, room.UpdatedAt())
	assert.Nil(t, room.deletedAt)

	room = NewRoomWith(id, adminId, name, category, createdAt, updateAt, deleteAt)
	assert.Equal(t, id.Value(), room.Id().Value())
	assert.Equal(t, adminId.Value(), room.AdminId().Value())
	assert.Equal(t, name.Value(), room.Name().Value())
	assert.Equal(t, category.Value(), room.Category().Value())
	assert.Equal(t, createdAt.Value(), room.CreatedAt().Value())
	assert.Equal(t, updateAt.Value(), room.UpdatedAt().Value())
	assert.Equal(t, deleteAt, room.DeletedAt())
}

func TestShouldUpdateARoomName(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room := NewRoom(adminId, name, category)

	oldUpdatedAt := room.UpdatedAt()
	newName, _ := valueobject.NewRoomNameWith("Rust")

	room.UpdateName(newName)
	assert.Equal(t, newName.Value(), room.Name().Value())
	assert.True(t, room.updatedAt.Time().After(oldUpdatedAt.Time()))
}

func TestShouldUpdateARoomCategory(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room := NewRoom(adminId, name, category)

	oldUpdatedAt := room.UpdatedAt()
	newCategory, _ := valueobject.NewRoomCategoryWith("Science")

	room.UpdateCategory(newCategory)
	assert.Equal(t, newCategory.Value(), room.Category().Value())
	assert.True(t, room.updatedAt.Time().After(oldUpdatedAt.Time()))
}

func TestShouldValidateARoomAdmin(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room := NewRoom(adminId, name, category)

	err := room.ValidateAdmin(adminId)
	assert.Nil(t, err)

	fakeAdminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533g")
	err = room.ValidateAdmin(fakeAdminId)
	assert.IsType(t, validation.UnauthorizedError(""), err)
	assert.Error(t, ErrInvalidRoomAdmin, err)
}

func TestShouldDeleteARoomn(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room := NewRoom(adminId, name, category)

	assert.Nil(t, room.DeletedAt())
	assert.False(t, room.IsDeleted())

	err := room.Delete()
	assert.Nil(t, err)

	deletedAt := room.DeletedAt()

	err = room.Delete()
	assert.IsType(t, validation.ValidationError(""), err)
	assert.Error(t, ErrRoomAlreadyDeleted, err)

	assert.True(t, deletedAt.Time().Equal(room.deletedAt.Time()))
}
