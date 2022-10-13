# go-sfomuseum-pmtiles

Opinionated SFO Museum package for working with Protomaps (v3) databases.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-sfomuseum-pmtiles.svg)](https://pkg.go.dev/github.com/sfomuseum/go-sfomuseum-pmtiles)

## Tools

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

```
$> ./bin/server -h
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
  -example-latitude string
    	The starting latitude for the example map application. (default "37.6143")
  -example-longitude string
    	The starting longitude for the example map application. (default "-122.3828")
  -example-zoom string
    	The starting zoom for the example map application. (default "13")
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
  -tile-path string
    	A valid gocloud.dev/blob bucket URI where .pmtiles databases are stored.
```

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

Or, with a PMTiles database hosted on S3:

```
$> bin/server \
	-tile-path 's3blob://{BUCKET}?region={REGION}&credentials={CREDENTIALS}' \
	-enable-example \
	-example-database \
	sfomuseum
```

![](docs/images/example-world.png)

#### AWS

##### Lambda (with API Gateway proxy integration)

```
$> make lambda
if test -f main; then rm -f main; fi
if test -f server.zip; then rm -f server.zip; fi
GOOS=linux go build -mod vendor -o main cmd/server/main.go
zip server.zip main
  adding: main (deflated 55%)
rm -f main
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

##### Lambda URL Functions

Not yet. It's on the list.

* https://aws.amazon.com/blogs/networking-and-content-delivery/using-amazon-cloudfront-with-aws-lambda-as-origin-to-accelerate-your-web-applications/
* https://github.com/aaronland/go-http-server/issues/14

##### API Gateway

* Create a new `{proxy+}` resource on "/".
* Add a new `GET` method (on the "/" resource) and point it to your Lambda function.
* Be sure to add an entry for "application/x-protobuf" in `API: {API_NAME} > Settings > Binary Media Types`.

## See also

* https://github.com/protomaps/go-pmtiles
* https://github.com/aaronland/go-http-server
* https://gocloud.dev/blob
* https://github.com/aaronland/gocloud-blob-s3
