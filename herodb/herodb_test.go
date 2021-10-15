package herodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_MakePersistent(t *testing.T) {
	db := New()
	mpersist := &MockPersist{}

	mpersist.On("Load", context.Background()).Return(map[string]string{"adi": "ali"}, nil).Once()
	mpersist.On("Save", context.Background(), map[string]string{"adi": "ali"}).Return(nil)

	err := db.MakePersistent(context.Background(), mpersist, time.Second*1)
	require.NoError(t, err)

	time.Sleep(time.Second * 2)

	mpersist.AssertExpectations(t)

}

func TestStore_Get(t *testing.T) {
	db := New()
	db.db = map[string]string{"adi": "ali"}

	t.Run("existing key", func(t *testing.T) {

		val, err := db.Get(context.Background(), "adi")

		assert.NoError(t, err)
		assert.Equal(t, "ali", val)
	})

	t.Run("not exists key", func(t *testing.T) {

		val, err := db.Get(context.Background(), "yas")

		assert.ErrorIs(t, err, ErrNotfound)
		assert.Equal(t, "", val)
	})
}
