// Package server implements a web application for serving Protomaps tiles.
package server

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	gohttp "net/http"

	"github.com/aaronland/go-http-rewrite"
	"github.com/aaronland/go-http-server"
	"github.com/protomaps/go-pmtiles/pmtiles"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/bucket"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/example"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/http"
)

type RunOptions struct {
	Logger               *log.Logger
	FS                   fs.FS
	EnableCORS           bool
	CORSOrigins          []string
	CORSAllowCredentials bool
	CORSDebug            bool
	HTTPServerURI        string

	// To do: Add pmtiles.Server vars here
	// To do: Add example_ vars here
}

func RunOptionsWithFlagSet(fs *flag.FlagSet, logger *log.Logger) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SFOMUSEUM")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		HTTPServerURI:        server_uri,
		Logger:               logger,
		EnableCORS:           enable_cors,
		CORSOrigins:          cors_origins,
		CORSAllowCredentials: cors_allow_credentials,
		CORSDebug:            cors_debug,
	}

	return opts, nil
}

// RunWithFlagSet runs the server application using a default flagset.
func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

// RunWithFlagSet runs the server application using flags derived from 'fs'.
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SFOMUSEUM")

	if err != nil {
		return fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		HTTPServerURI:        server_uri,
		Logger:               logger,
		EnableCORS:           enable_cors,
		CORSOrigins:          cors_origins,
		CORSAllowCredentials: cors_allow_credentials,
		CORSDebug:            cors_debug,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	var pmtiles_server *pmtiles.Server
	var err error

	if opts.FS != nil {

		fs_bucket, err := bucket.NewBucketWithFS(opts.FS, tile_path, "")

		if err != nil {
			return fmt.Errorf("Failed to create new bucket from filesystem, %w", err)
		}

		defer fs_bucket.Close()

		pmtiles_server, err = pmtiles.NewServerWithBucket(fs_bucket, "", opts.Logger, 64, "", "")

		if err != nil {
			return fmt.Errorf("Failed to create PMTiles server from bucket, %w", err)
		}

	} else {

		pmtiles_server, err = pmtiles.NewServer(tile_path, "", opts.Logger, 64, "", "")

		if err != nil {
			return fmt.Errorf("Failed to create PMTiles server, %w", err)
		}
	}

	pmtiles_server.Start()

	mux := gohttp.NewServeMux()

	tile_handler := http.TileHandler(pmtiles_server, opts.Logger)

	if opts.EnableCORS {

		if len(opts.CORSOrigins) == 0 {
			opts.CORSOrigins = []string{"*"}
		}

		c := cors.New(cors.Options{
			AllowedOrigins:   opts.CORSOrigins,
			AllowCredentials: opts.CORSAllowCredentials,
			Debug:            opts.CORSDebug,
		})

		tile_handler = c.Handler(tile_handler)
	}

	mux.Handle("/", tile_handler)

	if enable_example {

		if example_database == "" {
			return fmt.Errorf("You must specify a value for -example-database.")
		}

		append_opts := &rewrite.AppendResourcesOptions{
			DataAttributes: map[string]string{
				"example-database":  example_database,
				"example-latitude":  example_latitude,
				"example-longitude": example_longitude,
				"example-zoom":      example_zoom,
			},
		}

		http_fs := gohttp.FS(example.FS)
		example_handler := gohttp.FileServer(http_fs)

		example_handler = rewrite.AppendResourcesHandler(example_handler, append_opts)
		example_handler = gohttp.StripPrefix("/example", example_handler)

		mux.Handle("/example/", example_handler)
	}

	null_handler := http.NullHandler()
	mux.Handle("/favicon.ico", null_handler)

	s, err := server.NewServer(ctx, opts.HTTPServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	opts.Logger.Printf("Listening for requests on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}
