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

It is also possible to use a PMTiles database embedded in an `fs.FS` instance. This is functionality specific to the `go-sfomuseum-pmtiles` package
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

Which tells the code to copy the files in the embedded filesystem to a new [in-memory `gocloud.dev/blob.Bucket` instance](https://gocloud.dev/howto/blob/).

_Technically it copies the data to a locally-defined [`GoCloudBucket`](bucket/gocloud.go) instance that wraps the `blob.Bucket` instance and implements the `protomaps/go-pmtiles/pmtiles.Bucket` interface. Computers, amirite?_

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
| SFOMUSEUM_TILE_PATH | s3blob://{BUCKET}?prefix={PREFIX}&region={REGION}&credentials=iam: | Note the `s3blob://` scheme which is different that the default `s3://` scheme and supports specifying AWS credentials using the `?credentials` parameter. |

The rules for assigning flags from envinronment variables are:

* Replace all instances of "-" in a flag name with "_".
* Uppercase the flag name.
* Prepend the new string with "SFOMUSEUM_".

For example the `-tile-path` flag becomes the `SFOMUSEUM_TILE_PATH` environment variable.

You will also need to add the AWS web adapter layer:

```
arn:aws:lambda:{YOUR_AWS_REGION}:753240598075:layer:LambdaAdapterLayerArm64:20
```

You will need to configure your Lambda functions with an IAM role that allows the function to read data from the S3 bucket named `{BUCKET}`.

How you configure access to your Lambda function URL as well as any CORS details is left up to you.

For example:

```
$> curl -I https://{FUNCTION_URL_PREFIX}.lambda-url.us-east-1.on.aws/{PROTOMAPS_DATABASE_NAME}/12/655/1585.mvt
HTTP/1.1 200 OK
Date: Mon, 08 Apr 2024 19:06:52 GMT
Content-Type: application/x-protobuf
Content-Length: 0
Connection: keep-alive
x-amzn-RequestId: c87c2bf1-b21d-4ed2-9ea2-3c195455ac1a
content-encoding: gzip
etag: "9b1d497bc8170939"
X-Amzn-Trace-Id: root=1-6614404c-552ab05545ce154700958f21;parent=11cc270b2c395c7b;sampled=0;lineage=a7dcb3fd:0
```

##### Lambda Function URL (with CloudFront integration)

CloudFront integration requires a bit of hoop-jumping if you are trying to hang your Function URL from a "leaf" node (for example: `your-domain.com/pmtiles`) and you are serving tiles from a sub-folder in an S3 bucket.

The default configuration would be to add Lambda environment variables like this:

| Name | Value | Notes |
| --- | --- | --- |
| SFOMUSEUM_TILE_PATH | s3blob://{BUCKET}?prefix={PREFIX}&region={REGION}&credentials=iam: | Note the `s3blob://` scheme which is different that the default `s3://` scheme and supports specifying AWS credentials using the `?credentials` parameter. |

Instead what you'll need to do two things. First, configure your Lambda environment variables like this:

| Name | Value | Notes |
| --- | --- | --- |
| SFOMUSEUM_TILE_PATH | s3blob://{BUCKET}?region={REGION}&credentials=iam: | Note the `s3blob://` scheme which is different that the default `s3://` scheme and supports specifying AWS credentials using the `?credentials` parameter. |
| SFOMUSEUM_STRIP_PREFIX | {PREFIX} | |

Note the way `SFOMUSEUM_TILE_PATH` no longer has a `?prefix={PREFIX}` value.

And when you configure the CloudFront origin be sure to assign the "Origin path" value to be `{PREFIX}`. For example if `{PREFIX}` is "/pmtiles" then the actual raw Function URL address would be for any given tile would be:

```
https://{FUNCTION_URL_PREFIX}.lambda-url.us-east-1.on.aws/pmtiles/pmtiles/{PROTOMAPS_DATABASE_NAME}/12/655/1585.mvt
```

Here's what ends up happening:

* The first `/pmtiles` prefix get appended to the origin (Function URL) path
* The prefix is stripped by the `SFOMUSEUM_STRIP_PREFIX` directive
* That means the final path, used to derive PMTiles tiles, is `pmtiles/{PROTOMAPS_DATABASE_NAME}/{z}/{x}/{y}.mvt`

Honestly, I don't understand why this is necessary. It seems like it _should_ be possible to simply specify a CloudFront origin with no "path", a `SFOMUSEUM_TILE_PATH` value with no prefix and then simply rely on the path/URI being defined in the CloudFront behaviour to be passed to the Lambda Function URL.

For example `https://your-domain.com/pmtiles/{PROTOMAPS_DATABASE_NAME}` would map the `pmtiles/{PROTOMAPS_DATABASE_NAME}.pmtiles` file in the `s3blob://{BUCKET}` bucket. But apparently not, or maybe I am just "doing it wrong"?

Anyway, pending the "right way" to do it the steps away acheive the same result.

## See also

* https://github.com/protomaps/go-pmtiles
* https://github.com/aaronland/go-http-server
* https://gocloud.dev/blob
* https://github.com/aaronland/gocloud-blob-s3
