package blobstore

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiskBlobStore(t *testing.T) {
	for _, tc := range []struct {
		name             string
		prefix, blobName string
		blob             []byte
	}{
		{
			"WithoutPrefix",
			"",
			"test_blob",
			[]byte("test"),
		},
		{
			"WithPrefix",
			"my_prefix",
			"test_blob",
			[]byte("test"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			*pathPrefix = tc.prefix

			tmpDir := t.TempDir()
			dbs, err := NewDiskBlobStore(tmpDir)
			require.NoError(t, err)

			ctx := context.Background()
			n, err := dbs.WriteBlob(ctx, tc.blobName, tc.blob)
			require.NoError(t, err)
			require.Greater(t, n, 0)

			path := filepath.Join(tmpDir, tc.prefix, tc.blobName)
			require.FileExists(t, path)

			exist, err := dbs.BlobExists(ctx, tc.blobName)
			require.NoError(t, err)
			require.True(t, exist)

			b, err := dbs.ReadBlob(ctx, tc.blobName)
			require.NoError(t, err)
			require.Equal(t, b, tc.blob)

			err = dbs.DeleteBlob(ctx, tc.blobName)
			require.NoError(t, err)

			exist, err = dbs.BlobExists(ctx, tc.blobName)
			require.NoError(t, err)
			require.False(t, exist)
		})
	}
}
