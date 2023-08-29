LAST_TAG=$(shell git describe --abbrev=0 --tags)
CURR_SHA=$(shell git rev-parse --verify HEAD)

LDFLAGS=-ldflags "-s -w -X main.version=$(LAST_TAG)"

.PHONY: release

all: build

# make release tag=v0.4.3
release:
	git tag $(tag)
	git push origin $(tag)

build:
	GOOS=$(os) GOARCH=$(arch) go build ${LDFLAGS} -o bin/$(exe) ./cmd/gnols