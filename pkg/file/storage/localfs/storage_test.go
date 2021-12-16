// +build test_integration

package localfs

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/my-network/fsutil/pkg/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tests_my-network_fsutil_pkg_file_localfs")
	require.NoError(t, err)
	defer func() { assert.NoError(t, os.RemoveAll(tmpDir)) }()

	stor := NewStorage(tmpDir)
	ctx := context.Background()

	t.Run("file", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			obj, err := stor.Open(ctx, nil, file.Path{"file"}, file.FlagWrite|file.FlagCreate|file.FlagExcl, 0600)
			require.NoError(t, err)
			require.IsType(t, &File{}, obj)
		})
		t.Run("existing", func(t *testing.T) {
			obj, err := stor.Open(ctx, nil, file.Path{"file"}, file.FlagWrite, 0600)
			require.NoError(t, err)
			require.IsType(t, &File{}, obj)
		})
	})

	t.Run("dir", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			err := stor.Mkdir(ctx, nil, file.Path{"dir"}, 0700, false)
			require.NoError(t, err)
		})
	})

	t.Run("symlink", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			err := stor.Symlink(ctx, nil, file.Path{"dir", "symlink"}, file.Path{"..", "file"})
			require.NoError(t, err)
		})

		t.Run("open_nofollow", func(t *testing.T) {
			obj, err := stor.Open(ctx, nil, file.Path{"dir", "symlink"}, file.FlagWrite|file.FlagNoFollow|file.FlagPath, 0600)
			require.NoError(t, err)
			require.IsType(t, &Symlink{}, obj)
		})

		t.Run("open_follow", func(t *testing.T) {
			obj, err := stor.Open(ctx, nil, file.Path{"dir", "symlink"}, file.FlagWrite, 0600)
			require.NoError(t, err)
			require.IsType(t, &File{}, obj)
		})

		t.Run("open_abs", func(t *testing.T) {
			err := stor.Symlink(ctx, nil, file.Path{"dir", "symlink_abs"}, stor.ToAbsPath(file.Path{"file"}))
			require.NoError(t, err)

			obj, err := stor.Open(ctx, nil, file.Path{"dir", "symlink_abs"}, file.FlagWrite, 0600)
			require.NoError(t, err)
			require.IsType(t, &File{}, obj)
		})

		t.Run("readlink", func(t *testing.T) {
			destination, err := stor.Readlink(ctx, nil, file.Path{"dir", "symlink"})
			require.NoError(t, err)
			require.Equal(t, 2, len(destination), destination)
			require.Equal(t, "..", destination[0])
			require.Equal(t, "file", destination[1])
		})
	})
}
