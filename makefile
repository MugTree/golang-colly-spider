build:
	go build -o bin/spider ./...

run: build 
	./bin/spider

test: 
	go test -v ./... -count=1
