package herodb

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockPersist struct {
	mock.Mock
}

func (m *MockPersist) Load(ctx context.Context) (map[string]string, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockPersist) Save(ctx context.Context, data map[string]string) error {
	return m.Called(ctx, data).Error(0)
}
