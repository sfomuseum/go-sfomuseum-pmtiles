// Package server implements a web application for serving Protomaps tiles.
package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	gohttp "net/http"

	"github.com/aaronland/go-http-rewrite"
	"github.com/aaronland/go-http-server"
	"github.com/protomaps/go-pmtiles/pmtiles"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/example"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/http"
)

type RunOptions struct {
	Server               *pmtiles.Server
	Logger               *log.Logger
	EnableCORS           bool
	CORSOrigins          []string
	CORSAllowCredentials bool
	CORSDebug            bool
	HTTPServerURI        string
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

	server, err := pmtiles.NewServer(tile_path, "", logger, cache_size, "", "")

	if err != nil {
		return fmt.Errorf("Failed to create pmtiles.Loop, %w", err)
	}

	opts := &RunOptions{
		HTTPServerURI:        server_uri,
		Server:               server,
		Logger:               logger,
		EnableCORS:           enable_cors,
		CORSOrigins:          cors_origins,
		CORSAllowCredentials: cors_allow_credentials,
		CORSDebug:            cors_debug,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	opts.Server.Start()

	mux := gohttp.NewServeMux()

	tile_handler := http.TileHandler(opts.Server, opts.Logger)

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
