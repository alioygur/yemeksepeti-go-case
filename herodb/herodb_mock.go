package herodb

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (s *MockStore) Get(ctx context.Context, key string) (string, error) {
	args := s.Called(ctx, key)

	return args.String(0), args.Error(1)
}

func (s *MockStore) Set(ctx context.Context, key string, value string) error {
	return s.Called(ctx, key, value).Error(0)
}

func (s *MockStore) Flush(ctx context.Context) error {
	return s.Called(ctx).Error(0)
}
