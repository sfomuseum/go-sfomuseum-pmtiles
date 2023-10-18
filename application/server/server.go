// Package server implements a web application for serving Protomaps tiles.
package server

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	gohttp "net/http"
	"strconv"

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
	ServerURI string
	Logger    *log.Logger

	EnableCORS           bool
	CORSOrigins          []string
	CORSAllowCredentials bool
	CORSDebug            bool

	EnableExample    bool
	ExampleDatabase  string
	ExampleLatitude  float64
	ExampleLongitude float64
	ExampleZoom      int

	PMTilesURI         string
	PMTilesFS          fs.FS
	PMTilesPrefix      string
	PMTilesCacheSize   int
	PMTilesHostname    string
	PMTilesStripPrefix string
}

func RunOptionsWithFlagSet(fs *flag.FlagSet, logger *log.Logger) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SFOMUSEUM")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		ServerURI: server_uri,
		Logger:    logger,

		PMTilesURI:         tile_path,
		PMTilesPrefix:      tile_prefix,
		PMTilesCacheSize:   cache_size,
		PMTilesHostname:    public_hostname,
		PMTilesStripPrefix: strip_prefix,

		EnableCORS:           enable_cors,
		CORSOrigins:          cors_origins,
		CORSAllowCredentials: cors_allow_credentials,
		CORSDebug:            cors_debug,

		EnableExample:    enable_example,
		ExampleDatabase:  example_database,
		ExampleLatitude:  example_latitude,
		ExampleLongitude: example_longitude,
		ExampleZoom:      example_zoom,
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

	opts, err := RunOptionsWithFlagSet(fs, logger)

	if err != nil {
		return fmt.Errorf("Failed to derive run options, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	var pmtiles_server *pmtiles.Server
	var err error

	if opts.PMTilesFS != nil {

		fs_bucket, err := bucket.NewBucketWithFS(opts.PMTilesFS, opts.PMTilesURI, opts.PMTilesPrefix)

		if err != nil {
			return fmt.Errorf("Failed to create new bucket from filesystem, %w", err)
		}

		defer fs_bucket.Close()

		pmtiles_server, err = pmtiles.NewServerWithBucket(fs_bucket, opts.PMTilesPrefix, opts.Logger, opts.PMTilesCacheSize, "", opts.PMTilesHostname)

		if err != nil {
			return fmt.Errorf("Failed to create PMTiles server from bucket, %w", err)
		}

	} else {

		pmtiles_server, err = pmtiles.NewServer(tile_path, opts.PMTilesPrefix, opts.Logger, opts.PMTilesCacheSize, "", opts.PMTilesHostname)

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

	if opts.PMTilesStripPrefix != "" {
		tile_handler = gohttp.StripPrefix(opts.PMTilesStripPrefix, tile_handler)
	}

	mux.Handle("/", tile_handler)

	if opts.EnableExample {

		if opts.ExampleDatabase == "" {
			return fmt.Errorf("You must specify a value for -example-database.")
		}

		append_opts := &rewrite.AppendResourcesOptions{
			DataAttributes: map[string]string{
				"example-database":  opts.ExampleDatabase,
				"example-latitude":  strconv.FormatFloat(opts.ExampleLatitude, 'f', -1, 64),
				"example-longitude": strconv.FormatFloat(opts.ExampleLongitude, 'f', -1, 64),
				"example-zoom":      strconv.Itoa(opts.ExampleZoom),
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

	s, err := server.NewServer(ctx, opts.ServerURI)

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
