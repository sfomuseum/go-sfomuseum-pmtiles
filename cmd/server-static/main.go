// server launches a web server for search Protomaps (v3) tile requests.
package main

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob-s3"
	app "github.com/sfomuseum/go-sfomuseum-pmtiles/application/server"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/static"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	flag_fs := app.DefaultFlagSet()

	opts, err := app.RunOptionsWithFlagSetAndFS(flag_fs, logger, static.FS)

	if err != nil {
		logger.Fatalf("Failed to derive run options from flagset, %w", err)
	}

	err = app.RunWithOptions(ctx, opts)

	if err != nil {
		logger.Fatal("Failed to run server, %w", err)
	}
}
