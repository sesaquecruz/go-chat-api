package validation

import "errors"

var ErrRequiredId = errors.New("id is required")
var ErrInvalidId = errors.New("id is invalid")

var ErrRequiredRoomId = errors.New("room id is required")
var ErrRequiredRoomAdminId = errors.New("room admin id is required")
var ErrRequiredRoomName = errors.New("room name is required")
var ErrRequiredRoomCategory = errors.New("room category is required")
var ErrRequiredRoomCreatedAt = errors.New("room created at is required")
var ErrRequiredRoomUpdatedAt = errors.New("room updated at is required")
var ErrMaxSizeRoomName = errors.New("room name must not have more than 20 characters")
var ErrInvalidRoomCategory = errors.New("room category is invalid")
var ErrAlreadyExistsRoom = errors.New("room already exists")
var ErrNotFoundRoom = errors.New("room not found")
var ErrInvalidRoomAdmin = errors.New("invalid room admin")

var ErrRequiredTimestamp = errors.New("timestamp is required")
var ErrInvalidTimestamp = errors.New("timestamp is invalid")
