build:
	go mod download
	go build -o recmd-cli

default: build

test:
	(go test github.com/tarof429/recmd-cli/cli)

clean:
	rm -f recmd-cli

install:
	go install