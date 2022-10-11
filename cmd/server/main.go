package main

import (
	"context"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/application/server"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
	_ "github.com/aaronland/gocloud-blob-s3"
	"log"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		logger.Fatal("Failed to run server, %w", err)
	}
}
