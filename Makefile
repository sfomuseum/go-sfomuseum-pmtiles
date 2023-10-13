CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

debug:
	go run -mod $(GOMOD) cmd/server/main.go \
		-enable-example \
		-tile-path file://$(CWD)/fixtures \
		-example-database sfo

debug-static:
	go run -mod $(GOMOD) cmd/server-static/main.go \
		-enable-example \
		-example-database sfo \
		-tile-path mem://

cli:
	go build -ldflags="-s -w" -mod $(GOMOD) -o bin/server cmd/server/main.go 

lambda:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/server/main.go
	zip server.zip bootstrap
	rm -f bootstrap

lambda-static:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server-static.zip; then rm -f server-static.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/server-static/main.go
	zip server-static.zip bootstrap
	rm -f bootstrap
