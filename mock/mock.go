package mockWriteSyncer

import (
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockWriterSyncerExpect struct {
	mock *mock.Mock
}

func (ws *MockWriterSyncerExpect) WriteJson(t *testing.T, obj interface{}) *mock.Call {
	exp := ws.mock.On(
		"Write",
		mock.AnythingOfType("[]uint8"),
	)

	return exp.Run(func(args mock.Arguments) {
		data, ok := args.Get(0).([]byte)
		assert.True(t, ok)
		require.NoError(t, json.Unmarshal(data, obj))
		exp.Return(len(data), nil)
	})
}

func (ws *MockWriterSyncerExpect) Write(p []byte) *mock.Call {
	return ws.mock.On("Write", p)
}

func (ws *MockWriterSyncerExpect) Sync() *mock.Call {
	return ws.mock.On("Sync")
}

type MockWriteSyncer struct {
	expect *MockWriterSyncerExpect
}

func (ws *MockWriteSyncer) Write(p []byte) (n int, err error) {
	args := ws.expect.mock.Called(p)
	return args.Int(0), args.Error(1)
}

func (ws *MockWriteSyncer) Sync() error {
	args := ws.expect.mock.Called()
	return args.Error(0)
}

func (ws *MockWriteSyncer) EXPECT() *MockWriterSyncerExpect {
	return ws.expect
}

func New() *MockWriteSyncer {
	return &MockWriteSyncer{
		expect: &MockWriterSyncerExpect{
			mock: &mock.Mock{},
		},
	}
}
