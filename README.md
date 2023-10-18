# go-sfomuseum-pmtiles

Opinionated SFO Museum package for working with Protomaps (v3) databases.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-sfomuseum-pmtiles.svg)](https://pkg.go.dev/github.com/sfomuseum/go-sfomuseum-pmtiles)

## Tools

```
$> make cli
go build -ldflags="-s -w" -mod vendor -o bin/server cmd/server/main.go 
```

### server

```
$> ./bin/server -h
Launch a web server for search Protomaps (v3) tile requests.
Usage:
	 ./bin/server [options]
Valid options are:
  -cache-size int
    	Cache size in megabytes for tiles. (default 64)
  -cors-allow-credentials
    	Enable support for credentials in CORS requests.
  -cors-debug
    	Enable debugging in the rs/cors package.
  -cors-origin value
    	One or more comma-separated lists of hosts to enable CORS support for. If the -enable-cors flag is set and no -cors-origin flags have been assigned then CORS support will be enabled for '*'.
  -enable-cors
    	Enable CORS support.
  -enable-example
    	Enable an example map application at /example for testing database files. 
  -example-database string
    	The name of the database to use in the example map application. Note that this value should be the name of the database without its extension.
  -example-latitude float
    	The starting latitude for the example map application. (default 37.6143)
  -example-longitude float
    	The starting longitude for the example map application. (default -122.3828)
  -example-zoom int
    	The starting zoom for the example map application. (default 13)
  -public-hostname string
    	Public hostname of tile endpoint.
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
  -strip-prefix string
    	An optional string prefix to strip from HTTP request for tiles.
  -tile-path string
    	A valid gocloud.dev/blob bucket URI where .pmtiles databases are stored.
  -tile-prefix string
    	A prefix to append when fetching tiles.
```

### Examples

#### localhost

For example:

```
$> ./bin/server \
	-tile-path file:///usr/local/sfomuseum/tiles \
	-enable-example \
	-example-database sfo

2022/10/11 14:54:31 Listening for requests on http://localhost:8080
2022/10/11 14:54:37 fetching sfo 0-16384
2022/10/11 14:54:37 fetched sfo 0-0
2022/10/11 14:54:37 [200] served /sfo/11/328/792.mvt in 3.025132ms
2022/10/11 14:54:37 [200] served /sfo/11/327/792.mvt in 4.38898ms
2022/10/11 14:54:38 [200] served /sfo/12/656/1585.mvt in 322.892µs
2022/10/11 14:54:38 [200] served /sfo/12/654/1585.mvt in 7.676399ms
2022/10/11 14:54:38 [200] served /sfo/12/655/1585.mvt in 7.976492ms
2022/10/11 14:54:39 [200] served /sfo/13/1311/3170.mvt in 9.82361ms
...and so on
```

![](docs/images/example-sfo.png)

#### S3

Or, with a PMTiles database hosted on S3:

```
$> bin/server \
	-tile-path 's3blob://{BUCKET}?region={REGION}&credentials={CREDENTIALS}' \
	-enable-example \
	-example-database \
	sfomuseum
```

![](docs/images/example-world.png)

#### static/embedded

It is also possible to use an PMTiles database embedded in an `fs.FS` instance. This is functionality specific to the `go-sfomuseum-pmtiles` package
rather than `protomaps/go-pmtiles` itself.

In order to use an embedded PMTiles database you need to explictly define a `fs.FS` instance where databases are stored. The contents of the `fs.FS` filesystem are copied to a `gocloud.dev/blob.Bucket` at runtime which is then used to initialize a `go-pmtiles.Bucket` instance for serving tiles.

For example:

```
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

	opts, _ := app.RunOptionsWithFlagSet(flag_fs, logger)
	opts.PMTilesFS = static.FS

	app.RunWithOptions(ctx, opts)
}
```

_Error handling omitted for the sake of brevity._

This functionality is not suitable for very large PMTiles databases but might be useful for applications with small, bounded, maps. To view a working example, run the `debug-static` Makefile target and then open a web browser at `http://localhost:8080/example`.

```
$> make debug-static
go run -mod readonly cmd/server-static/main.go \
		-enable-example \
		-example-database sfo \
		-tile-path mem://
		
2023/10/13 16:25:33 Listening for requests on http://localhost:8080
```

_Screenhost omitted because it looks the same as the screenshot for local hosting._

The relevant bits here are this:

```
	opts, _ := app.RunOptionsWithFlagSet(flag_fs, logger)
	opts.PMTilesFS = static.FS
```

And this:

```
	-example-database sfo
```

Which tells the code to load a PMTiles database named "sfo.db" from an [embedded filesystem](https://pkg.go.dev/embed) defined in [static/static.go](static/static.go). The other relevant flag is:

```
	-tile-path mem://
```

Which tells the code to copy the files in the embedded filesystem to a new [in-memory `gocloud.dev/blob.Bucket` instance](https://gocloud.dev/howto/blob/). Technically it copies the data to a locally-defined [`GoCloudBucket`](bucket/gocloud.go) instance that wraps the `blob.Bucket` instance and implements the `protomaps/go-pmtiles/pmtiles.Bucket` interface. Computers, amirite?

If instead you wanted to copy the embedded filesystem to a local filesystem you would update the `-tile-path` parameter to specify the relevant "file:///path/to/folder" URI. This might seem counter-intuitive but if used in a frequently-invoked AWS Lambda function context you could take advantage of the fact the function (and its filesystem) will persist across invocations and bundle a PMTiles database that might otherwise be too large to bundle in memory without the need to setup and configure an S3 bucket to store the database(s) you are serving.

The point is not that you _should_ do it this way. The point is that there are circumstances where you might want or need to serve tiles from an embedded filesystem and now you can.

By default the `server-static` tool supports cloning embedded filesystems to memory or a local filesystem. If you want to copy filesystems to other [blob.Bucket](https://gocloud.dev/howto/blob/) implementations you will need to [clone the code](cmd/server-static/main.go) and add the relevent `import` statements.

#### AWS

##### Lambda (with Function URL integration)

```
$> make lambda
if test -f bootstrap; then rm -f bootstrap; fi
if test -f server.zip; then rm -f server.zip; fi
GOARCH=arm64 GOOS=linux go build -mod readonly -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/server/main.go
zip server.zip bootstrap
  adding: bootstrap (deflated 71%)
rm -f bootstrap
```

The following environment variables should be configured for use as a Lambda function:

| Name | Value | Notes |
| --- | --- | --- |
| SFOMUSEUM_SERVER_URI | functionurl://?binary_type=application/x-protobuf | Note the `?binary_type` parameter. This is important. |
| SFOMUSEUM_TILE_PATH | s3blob://{BUCKET}?prefix={PREFIX}&region={REGION}&credentials=iam: | Note the `s3blob://` scheme which is different that the default `s3://` scheme and supports specifying AWS credentials using the `?credentials` parameter. |

The rules for assigning flags from envinronment variables are:

* Replace all instances of "-" in a flag name with "_".
* Uppercase the flag name.
* Prepend the new string with "SFOMUSEUM_".

For example the `-server-uri` flag becomes the `SFOMUSEUM_SERVER_URI` environment variable.

You will need to configure your Lambda functions with an IAM role that allows the function to read data from the S3 bucket named `{BUCKET}`.

How you configure access to your Lambda function URL as well as any CORS details is left up to you.

##### Lambda (with API Gateway proxy integration)

```
$> make lambda
if test -f bootstrap; then rm -f bootstrap; fi
if test -f server.zip; then rm -f server.zip; fi
GOARCH=arm64 GOOS=linux go build -mod readonly -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/server/main.go
zip server.zip bootstrap
  adding: bootstrap (deflated 71%)
rm -f bootstrap
```

The following environment variables should be configured for use as a Lambda function:

| Name | Value | Notes |
| --- | --- | --- |
| SFOMUSEUM_SERVER_URI | lambda://?binary_type=application/x-protobuf | Note the `?binary_type` parameter. This is important. |
| SFOMUSEUM_TILE_PATH | s3blob://{BUCKET}?prefix={PREFIX}&region={REGION}&credentials=iam: | Note the `s3blob://` scheme which is different that the default `s3://` scheme and supports specifying AWS credentials using the `?credentials` parameter. |
| SFOMUSEUM_CORS_ENABLE | true | |

The rules for assigning flags from envinronment variables are:

* Replace all instances of "-" in a flag name with "_".
* Uppercase the flag name.
* Prepend the new string with "SFOMUSEUM_".

For example the `-server-uri` flag becomes the `SFOMUSEUM_SERVER_URI` environment variable.

You will need to configure your Lambda functions with an IAM role that allows the function to read data from the S3 bucket named `{BUCKET}`.

##### API Gateway

* Create a new `{proxy+}` resource on "/".
* Add a new `GET` method (on the "/" resource) and point it to your Lambda function.
* Be sure to add an entry for "application/x-protobuf" in `API: {API_NAME} > Settings > Binary Media Types`.

## See also

* https://github.com/protomaps/go-pmtiles
* https://github.com/aaronland/go-http-server
* https://gocloud.dev/blob
* https://github.com/aaronland/gocloud-blob-s3
