// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	sql "database/sql"
	reflect "reflect"
	time "time"

	domain "github.com/ali-aidaruly/outbox/internal/domain"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// Mockrepository is a mock of repository interface.
type Mockrepository struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryMockRecorder
}

// MockrepositoryMockRecorder is the mock recorder for Mockrepository.
type MockrepositoryMockRecorder struct {
	mock *Mockrepository
}

// NewMockrepository creates a new mock instance.
func NewMockrepository(ctrl *gomock.Controller) *Mockrepository {
	mock := &Mockrepository{ctrl: ctrl}
	mock.recorder = &MockrepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepository) EXPECT() *MockrepositoryMockRecorder {
	return m.recorder
}

// CreateMessage mocks base method.
func (m *Mockrepository) CreateMessage(ctx context.Context, tx *sql.Tx, msgType string, payload []byte) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", ctx, tx, msgType, payload)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMessage indicates an expected call of CreateMessage.
func (mr *MockrepositoryMockRecorder) CreateMessage(ctx, tx, msgType, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*Mockrepository)(nil).CreateMessage), ctx, tx, msgType, payload)
}

// FetchNextBatch mocks base method.
func (m *Mockrepository) FetchNextBatch(ctx context.Context, batchSize int, lease time.Duration) ([]domain.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchNextBatch", ctx, batchSize, lease)
	ret0, _ := ret[0].([]domain.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchNextBatch indicates an expected call of FetchNextBatch.
func (mr *MockrepositoryMockRecorder) FetchNextBatch(ctx, batchSize, lease interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchNextBatch", reflect.TypeOf((*Mockrepository)(nil).FetchNextBatch), ctx, batchSize, lease)
}

// MarkMessagesProcessed mocks base method.
func (m *Mockrepository) MarkMessagesProcessed(ctx context.Context, ids []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkMessagesProcessed", ctx, ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkMessagesProcessed indicates an expected call of MarkMessagesProcessed.
func (mr *MockrepositoryMockRecorder) MarkMessagesProcessed(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessagesProcessed", reflect.TypeOf((*Mockrepository)(nil).MarkMessagesProcessed), ctx, ids)
}
