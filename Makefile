build:
	go mod download
	go build -o recmd-cli

default: build

test:
	(cd recmd; go test)

clean:
	rm -f recmd

install:
	go install