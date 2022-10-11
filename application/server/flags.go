package server

import (
	"flag"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var server_uri string
var tile_path string
var cache_size int
var enable_cors bool
var cors_allow_credentials bool
var cors_debug bool

var cors_origins multi.MultiCSVString


var example bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("pmtiles")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&tile_path, "tile-path", "", "...")
	fs.IntVar(&cache_size, "cache-size", 64, "...")
	fs.BoolVar(&enable_cors, "enable-cors", false, "...")
	fs.BoolVar(&cors_allow_credentials, "cors-allow-credentials", false, "...")
	fs.BoolVar(&cors_debug, "cors-debug", false, "...")		

	fs.Var(&cors_origins, "cors-origin", "")

	fs.BoolVar(&example, "example", false, "...")
	return fs
}
