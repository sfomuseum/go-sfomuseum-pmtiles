package pmtiles

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"

	pm "github.com/protomaps/go-pmtiles/pmtiles"
)

func NewServerWithFS(bucket_fs fs.FS, bucketURL string, prefix string, logger *log.Logger, cacheSize int, cors string, publicHostname string) (*pm.Server, error) {

	ctx := context.Background()

	bucketURL, _, err := pm.NormalizeBucketKey(bucketURL, prefix, "")

	if err != nil {
		return nil, err
	}

	gc_bucket, err := NewGoCloudBucket(ctx, bucketURL, prefix)

	if err != nil {
		return nil, fmt.Errorf("Failed to open bucket, %v", err)
	}

	var walk_func func(path string, d fs.DirEntry, err error) error

	walk_func = func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return fmt.Errorf("Failed to walk %s, %w", path, err)
		}

		if d.IsDir() {

			if path == "." {
				return nil
			}

			return fs.WalkDir(bucket_fs, path, walk_func)
		}

		r, err := bucket_fs.Open(path)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		wr, err := gc_bucket.NewWriter(ctx, path, nil)

		if err != nil {
			return fmt.Errorf("Failed to create %s for writing, %w", path, err)
		}

		_, err = io.Copy(wr, r)

		if err != nil {
			return fmt.Errorf("Failed to copy %s, %w", path, err)
		}

		err = wr.Close()

		if err != nil {
			return fmt.Errorf("Failed to close %s, %w", path, err)
		}

		return nil
	}

	err = fs.WalkDir(bucket_fs, ".", walk_func)

	if err != nil {
		return nil, fmt.Errorf("Failed to walk filesystem, %w", err)
	}

	return pm.NewServerWithBucket(gc_bucket, prefix, logger, cacheSize, cors, publicHostname)
}
