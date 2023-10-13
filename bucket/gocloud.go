package bucket

import (
	"context"
	"fmt"
	"io"

	"github.com/protomaps/go-pmtiles/pmtiles"
	"gocloud.dev/blob"
)

type GoCloudBucket struct {
	pmtiles.Bucket
	bucket *blob.Bucket
}

func NewGoCloudBucket(ctx context.Context, uri string, prefix string) (*GoCloudBucket, error) {

	b, err := blob.OpenBucket(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open bucket, %w", err)
	}

	if prefix != "" {
		b = blob.PrefixedBucket(b, prefix)
	}

	gcb := &GoCloudBucket{
		bucket: b,
	}

	return gcb, nil
}

func (gcb *GoCloudBucket) Close() error {
	return gcb.bucket.Close()
}

func (gcb *GoCloudBucket) NewRangeReader(ctx context.Context, key string, offset, length int64) (io.ReadCloser, error) {
	return gcb.bucket.NewRangeReader(ctx, key, offset, length, nil)
}

func (gcb *GoCloudBucket) NewWriter(ctx context.Context, key string, opts *blob.WriterOptions) (*blob.Writer, error) {
	return gcb.bucket.NewWriter(ctx, key, opts)
}
