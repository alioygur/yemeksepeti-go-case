package herodb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var filename = "./testdata/data.db"

func TestNewFilepersist(t *testing.T) {
	t.Run("existsting db", func(t *testing.T) {
		_, err := NewFilepersist(filename)
		assert.NoError(t, err)
	})
	t.Run("not exists db", func(t *testing.T) {
		file := "./testdata/notexists.db"
		defer os.Remove(file)
		_, err := NewFilepersist(file)
		assert.NoError(t, err)
		assert.FileExists(t, file)
	})
}

func TestFilePersist_Load(t *testing.T) {
	f, err := NewFilepersist(filename)
	require.NoError(t, err)

	data, err := f.Load(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, make(map[string]string), data)
}

func TestFilePersist_Save(t *testing.T) {
	file := "./testdata/temp.db"
	f, err := NewFilepersist(file)
	require.NoError(t, err)
	defer os.Remove(file)

	err = f.Save(context.Background(), map[string]string{"name": "ali"})
	require.NoError(t, err)

	data, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("%s\n", `{"name":"ali"}`), string(data))
}
