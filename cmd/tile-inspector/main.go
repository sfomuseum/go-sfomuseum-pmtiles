package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/paulmach/orb/encoding/mvt"
)

func main() {

	var gzipped bool

	flag.BoolVar(&gzipped, "gz", false, "...")
	flag.Parse()

	for _, path := range flag.Args() {

		var r io.Reader

		fl_r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer fl_r.Close()

		if gzipped {

			gz_r, err := gzip.NewReader(fl_r)

			if err != nil {
				log.Fatalf("Failed to create gzip reader for %s, %v", path, err)
			}

			r = gz_r
		} else {
			r = fl_r
		}

		body, err := io.ReadAll(r)

		if err != nil {
			log.Fatalf("Failed to read %s, %v", path, err)
		}

		layers, err := mvt.Unmarshal(body)

		if err != nil {
			log.Fatalf("Failed to unmarshal %s, %v", path, err)
		}

		fc := layers.ToFeatureCollections()

		enc := json.NewEncoder(os.Stdout)
		err = enc.Encode(fc)

		if err != nil {
			log.Fatalf("Failed to encode layers for %s, %v", path, err)
		}
	}
}
