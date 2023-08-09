// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/gateway/room_gateway.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	valueobject "github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	gomock "go.uber.org/mock/gomock"
)

// MockRoomGatewayInterface is a mock of RoomGatewayInterface interface.
type MockRoomGatewayInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRoomGatewayInterfaceMockRecorder
}

// MockRoomGatewayInterfaceMockRecorder is the mock recorder for MockRoomGatewayInterface.
type MockRoomGatewayInterfaceMockRecorder struct {
	mock *MockRoomGatewayInterface
}

// NewMockRoomGatewayInterface creates a new mock instance.
func NewMockRoomGatewayInterface(ctrl *gomock.Controller) *MockRoomGatewayInterface {
	mock := &MockRoomGatewayInterface{ctrl: ctrl}
	mock.recorder = &MockRoomGatewayInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoomGatewayInterface) EXPECT() *MockRoomGatewayInterfaceMockRecorder {
	return m.recorder
}

// FindById mocks base method.
func (m *MockRoomGatewayInterface) FindById(ctx context.Context, id *valueobject.ID) (*entity.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, id)
	ret0, _ := ret[0].(*entity.Room)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockRoomGatewayInterfaceMockRecorder) FindById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockRoomGatewayInterface)(nil).FindById), ctx, id)
}

// Save mocks base method.
func (m *MockRoomGatewayInterface) Save(ctx context.Context, room *entity.Room) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, room)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockRoomGatewayInterfaceMockRecorder) Save(ctx, room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRoomGatewayInterface)(nil).Save), ctx, room)
}