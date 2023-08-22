package validation

import "errors"

// Id errors
var ErrRequiredId = errors.New("id is required")
var ErrInvalidId = errors.New("id is invalid")

// User errors
var ErrRequiredUserId = errors.New("user id is required")
var ErrInvalidUserId = errors.New("user id is invalid")

var ErrRequiredUserName = errors.New("user name is required")
var ErrSizeUserName = errors.New("user name must not have more than 50 characters")

// Room errors
var ErrRequiredRoomId = errors.New("room id is required")
var ErrRequiredRoomAdminId = errors.New("room admin id is required")
var ErrRequiredRoomName = errors.New("room name is required")
var ErrRequiredRoomCategory = errors.New("room category is required")
var ErrRequiredRoomCreatedAt = errors.New("room created at is required")
var ErrRequiredRoomUpdatedAt = errors.New("room updated at is required")

var ErrInvalidRoomAdmin = errors.New("room admin is invalid")
var ErrSizeRoomName = errors.New("room name must not have more than 50 characters")
var ErrInvalidRoomCategory = errors.New("room category is invalid")

var ErrAlreadyExistsRoom = errors.New("room already exists")
var ErrNotFoundRoom = errors.New("room not found")

// Timestamp errors
var ErrRequiredTimestamp = errors.New("timestamp is required")
var ErrInvalidTimestamp = errors.New("timestamp is invalid")

// Query errors
var ErrInvalidQueryPage = errors.New("query 'page' must be greater than or equal to 0")
var ErrInvalidQuerySize = errors.New("query 'size' must be greater than 1 and less than or equal to 50")
var ErrInvalidQuerySort = errors.New("query 'sort' must be 'asc' or 'desc'")
var ErrMaxQuerySearch = errors.New("query 'search' length must be less than or equal to 50")

// Message errors
var ErrRequiredMessageId = errors.New("message id is required")
var ErrRequiredMessageRoomId = errors.New("message room id is required")
var ErrRequiredMessageSenderId = errors.New("message sender id is required")
var ErrRequiredMessageSenderName = errors.New("message sender name is required")
var ErrRequiredMessageText = errors.New("message text is required")
var ErrRequiredMessageCreatedAt = errors.New("message created at is required")

var ErrSizeMessageText = errors.New("message test length must be less than or equal to 100")
