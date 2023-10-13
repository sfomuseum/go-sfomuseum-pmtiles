package server

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

// A valid aaronland/go-http-server URI.
var server_uri string

// A valid gocloud.dev/blob bucket URI where .pmtiles databases are stored.
var tile_path string

// Cache size in megabytes for tiles.
var cache_size int

// A prefix to append when fetching tiles.
var tile_prefix string

// Public hostname of tile endpoint
var public_hostname string

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
var example_latitude float64

// The starting longitude for the example map application.
var example_longitude float64

// The starting zoom for the example map application.
var example_zoom int

// DefaultFlagSet returns a `flag.FlagSet` instance configured with the default flags necessary for the server application.
func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("pmtiles")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")
	fs.StringVar(&tile_path, "tile-path", "", "A valid gocloud.dev/blob bucket URI where .pmtiles databases are stored.")
	fs.IntVar(&cache_size, "cache-size", 64, "Cache size in megabytes for tiles.")
	fs.StringVar(&tile_prefix, "tile-prefix", "", "A prefix to append when fetching tiles.")
	fs.StringVar(&public_hostname, "public-hostname", "", "Public hostname of tile endpoint.")
	fs.BoolVar(&enable_cors, "enable-cors", false, "Enable CORS support.")
	fs.BoolVar(&cors_allow_credentials, "cors-allow-credentials", false, "Enable support for credentials in CORS requests.")
	fs.BoolVar(&cors_debug, "cors-debug", false, "Enable debugging in the rs/cors package.")

	fs.Var(&cors_origins, "cors-origin", "One or more comma-separated lists of hosts to enable CORS support for. If the -enable-cors flag is set and no -cors-origin flags have been assigned then CORS support will be enabled for '*'.")

	fs.BoolVar(&enable_example, "enable-example", false, "Enable an example map application at /example for testing database files. ")

	fs.StringVar(&example_database, "example-database", "", "The name of the database to use in the example map application. Note that this value should be the name of the database without its extension.")

	fs.Float64Var(&example_latitude, "example-latitude", 37.6143, "The starting latitude for the example map application.")
	fs.Float64Var(&example_longitude, "example-longitude", -122.3828, "The starting longitude for the example map application.")
	fs.IntVar(&example_zoom, "example-zoom", 13, "The starting zoom for the example map application.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Launch a web server for search Protomaps (v3) tile requests.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
