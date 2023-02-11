CWD=$(shell pwd)

debug:
	go run -mod vendor cmd/server/main.go \
		-enable-example \
		-tile-path file://$(CWD)/fixtures \
		-example-database sfo

cli:
	go build -ldflags="-s -w" -mod vendor -o bin/server cmd/server/main.go 

lambda:
	if test -f main; then rm -f main; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux go build -ldflags="-s -w" -mod vendor -o main cmd/server/main.go
	zip server.zip main
	rm -f main
