// server-static launches a web server for serving Protomaps (v3) tile requests using an embedded Protomaps database.
package main

import (
	_ "gocloud.dev/blob/fileblob"	
	_ "gocloud.dev/blob/memblob"
)

import (
	"context"
	"log"

	app "github.com/sfomuseum/go-sfomuseum-pmtiles/application/server"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/static"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	flag_fs := app.DefaultFlagSet()

	opts, err := app.RunOptionsWithFlagSet(flag_fs, logger)

	if err != nil {
		logger.Fatalf("Failed to derive run options from flagset, %w", err)
	}

	opts.PMTilesFS = static.FS

	err = app.RunWithOptions(ctx, opts)

	if err != nil {
		logger.Fatal("Failed to run server, %w", err)
	}
}
