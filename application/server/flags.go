package server

import (
	"flag"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

// A valid aaronland/go-http-server URI.
var server_uri string

// A valid gocloud.dev/blob bucket URI where .pmtiles databases are stored.
var tile_path string

// Cache size in megabytes for tiles.
var cache_size int

// Enable CORS support for HTTP requests.
var enable_cors bool

// Enable support for credentials in CORS requests.
var cors_allow_credentials bool

// Enable debugging in the rs/cors package.
var cors_debug bool

// One or more comma-separated lists of hosts to enable CORS support for. If the -enable-cors flag is set and no -cors-origin flags have been assigned then CORS support will be enabled for '*'.
var cors_origins multi.MultiCSVString

// Enable an example web application at /example for testing database files.
var enable_example bool

// The name of the database to use in the example web application. Note that this value should be the name of the database without its extension.
var example_database string

// The starting latitude for the example map application.
var example_latitude string

// The starting longitude for the example map application.
var example_longitude string

// The starting zoom for the example map application.
var example_zoom string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("pmtiles")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")
	fs.StringVar(&tile_path, "tile-path", "", "A valid gocloud.dev/blob bucket URI where .pmtiles databases are stored.")
	fs.IntVar(&cache_size, "cache-size", 64, "Cache size in megabytes for tiles.")
	fs.BoolVar(&enable_cors, "enable-cors", false, "Enable CORS support.")
	fs.BoolVar(&cors_allow_credentials, "cors-allow-credentials", false, "Enable support for credentials in CORS requests.")
	fs.BoolVar(&cors_debug, "cors-debug", false, "Enable debugging in the rs/cors package.")

	fs.Var(&cors_origins, "cors-origin", "One or more comma-separated lists of hosts to enable CORS support for. If the -enable-cors flag is set and no -cors-origin flags have been assigned then CORS support will be enabled for '*'.")

	fs.BoolVar(&enable_example, "enable-example", false, "Enable an example map application at /example for testing database files. ")

	fs.StringVar(&example_database, "example-database", "", "The name of the database to use in the example map application. Note that this value should be the name of the database without its extension.")

	fs.StringVar(&example_latitude, "example-latitude", "37.6143", "The starting latitude for the example map application.")
	fs.StringVar(&example_longitude, "example-longitude", "-122.3828", "The starting longitude for the example map application.")
	fs.StringVar(&example_zoom, "example-zoom", "13", "The starting zoom for the example map application.")

	return fs
}
