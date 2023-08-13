package entity

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
)

func TestRoom_ShouldCreateANewRoomWhenValuesAreNotNil(t *testing.T) {
	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	createdAt := valueobject.NewDateTime()
	updateAt := valueobject.NewDateTime()

	room, err := NewRoom(adminId, name, category)
	assert.NotNil(t, room)
	assert.Nil(t, err)
	assert.Nil(t, room.Validate())
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
	assert.Nil(t, room.Validate())
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
	assert.Equal(t, createdAt.StringValue(), room.CreatedAt().StringValue())
	assert.Equal(t, updateAt.StringValue(), room.UpdatedAt().StringValue())
}

func TestRoom_ShouldReturnARequireIdErrorWhenIdIsNil(t *testing.T) {
	room, err := NewRoomWith(nil, nil, nil, nil, nil, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomId)
}

func TestRoom_ShouldReturnARequireAdminIdErrorWhenAdminIdIsNil(t *testing.T) {
	id := valueobject.NewID()

	room, err := NewRoom(nil, nil, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomAdminId)

	room, err = NewRoomWith(id, nil, nil, nil, nil, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomAdminId)
}

func TestRoom_ShouldReturnARequireNameErrorWhenNameIsNil(t *testing.T) {
	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")

	room, err := NewRoom(adminId, nil, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomName)

	room, err = NewRoomWith(id, adminId, nil, nil, nil, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomName)
}

func TestRoom_ShouldReturnARequireCategoryErrorWhenCategoryIsNil(t *testing.T) {
	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")

	room, err := NewRoom(adminId, name, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomCategory)

	room, err = NewRoomWith(id, adminId, name, nil, nil, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomCategory)
}

func TestRoom_ShouldReturnARequireCreatedAtErrorWhenCreatedAtIsNil(t *testing.T) {
	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")

	room, err := NewRoomWith(id, adminId, name, category, nil, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomCreatedAt)
}

func TestRoom_ShouldReturnARequireUpdatedAtErrorWhenUpdatedAtIsNil(t *testing.T) {
	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	createdAt := valueobject.NewDateTime()

	room, err := NewRoomWith(id, adminId, name, category, createdAt, nil)
	assert.Nil(t, room)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomUpdatedAt)
}

func TestRoom_ShouldUpdateRoomNameWhenNameIsNotNil(t *testing.T) {
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")

	room, err := NewRoom(adminId, name, category)
	assert.NotNil(t, room)
	assert.Nil(t, err)

	updatedAt := room.UpdatedAt()
	newName, _ := valueobject.NewRoomNameWith("Rust")

	err = room.UpdateName(newName)
	assert.Nil(t, err)
	assert.Equal(t, newName.Value(), room.Name().Value())
	assert.True(t, room.updatedAt.TimeValue().After(*updatedAt.TimeValue()))
}

func TestRoom_ShouldReturnARequiredNameErrorWhenNewNameNil(t *testing.T) {
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")

	room, err := NewRoom(adminId, name, category)
	assert.NotNil(t, room)
	assert.Nil(t, err)

	updatedAt := room.UpdatedAt()
	oldName := room.Name()

	err = room.UpdateName(nil)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomName)
	assert.Equal(t, oldName.Value(), room.Name().Value())
	assert.True(t, room.updatedAt.TimeValue().Equal(*updatedAt.TimeValue()))
}

func TestRoom_ShouldUpdateRoomCategoryWhenCategoryIsNotNil(t *testing.T) {
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")

	room, err := NewRoom(adminId, name, category)
	assert.NotNil(t, room)
	assert.Nil(t, err)

	updatedAt := room.UpdatedAt()
	newCategory, _ := valueobject.NewRoomCategoryWith("Science")

	err = room.UpdateCategory(newCategory)
	assert.Nil(t, err)
	assert.Equal(t, newCategory.Value(), room.Category().Value())
	assert.True(t, room.updatedAt.TimeValue().After(*updatedAt.TimeValue()))
}

func TestRoom_ShouldReturnARequiredCategoryErrorWhenNewCategoryNil(t *testing.T) {
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Golang")
	category, _ := valueobject.NewRoomCategoryWith("Tech")

	room, err := NewRoom(adminId, name, category)
	assert.NotNil(t, room)
	assert.Nil(t, err)

	updatedAt := room.UpdatedAt()
	oldCategory := room.Category()

	err = room.UpdateCategory(nil)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomCategory)
	assert.Equal(t, oldCategory.Value(), room.Category().Value())
	assert.True(t, room.updatedAt.TimeValue().Equal(*updatedAt.TimeValue()))
}
