package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-rewrite"
	"github.com/aaronland/go-http-server"
	"github.com/protomaps/go-pmtiles/pmtiles"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/example/www"
	"github.com/sfomuseum/go-sfomuseum-pmtiles/http"
	"log"
	gohttp "net/http"
)

// RunWithFlagSet runs the server application using a default flagset.
func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

// RunWithFlagSet runs the server application using flags derived from 'fs'.
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	loop := pmtiles.NewLoop(tile_path, logger, cache_size, "")
	loop.Start()

	mux := gohttp.NewServeMux()

	tile_handler := http.TileHandler(loop, logger)

	if enable_cors {

		if len(cors_origins) == 0 {
			cors_origins.Set("*")
		}

		c := cors.New(cors.Options{
			AllowedOrigins:   cors_origins,
			AllowCredentials: cors_allow_credentials,
			Debug:            cors_debug,
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

		http_fs := gohttp.FS(www.FS)
		example_handler := gohttp.FileServer(http_fs)

		example_handler = rewrite.AppendResourcesHandler(example_handler, append_opts)
		example_handler = gohttp.StripPrefix("/example", example_handler)

		mux.Handle("/example/", example_handler)
	}

	null_handler := http.NullHandler()
	mux.Handle("/favicon.ico", null_handler)

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	logger.Printf("Listening for requests on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}
