build:
	go mod download
	go build -o recmd

default: build

test:
	go test

clean:
	rm -f recmd