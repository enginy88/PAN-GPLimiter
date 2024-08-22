CC=go build
LDFLAGS="-s -w"
NAME=$(shell go list -m)

.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p ./bin
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 $(CC) -ldflags $(LDFLAGS) -o bin/$(NAME)_mac-amd64
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 $(CC) -ldflags $(LDFLAGS) -o bin/$(NAME)_mac-arm64
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 $(CC) -ldflags $(LDFLAGS) -o bin/$(NAME)_lin-amd64
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 $(CC) -ldflags $(LDFLAGS) -o bin/$(NAME)_lin-arm64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(CC) -ldflags $(LDFLAGS) -o bin/$(NAME)_win-amd64.exe
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 $(CC) -ldflags $(LDFLAGS) -o bin/$(NAME)_win-arm64.exe

.PHONY: local
local:
	mkdir -p ./bin
	CGO_ENABLED=0 $(CC) -ldflags $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf ./bin
