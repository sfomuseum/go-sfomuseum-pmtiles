// server launches a web server for search Protomaps (v3) tile requests.
package main

import (
	"context"
	_ "github.com/aaronland/gocloud-blob-s3"
	app "github.com/sfomuseum/go-sfomuseum-pmtiles/application/server"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
	"log"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := app.Run(ctx, logger)

	if err != nil {
		logger.Fatal("Failed to run server, %w", err)
	}
}
