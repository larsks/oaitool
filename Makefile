PKG=github.com/larsks/oaitool
EXE=oaitool-$(shell go env GOOS)-$(shell go env GOARCH)

GOSRC =  main.go \
	 $(wildcard api/*.go) \
	 $(wildcard cli/*.go) \
	 $(wildcard version/*.go)

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo unknown)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

GOLDFLAGS = \
	    -X '$(PKG)/version.BuildVersion=$(VERSION)' \
	    -X '$(PKG)/version.BuildRef=$(COMMIT)' \
	    -X '$(PKG)/version.BuildDate=$(DATE)'

all: $(EXE)

$(EXE): $(GOSRC)
	go build -o $@ -ldflags "$(GOLDFLAGS)"

clean:
	rm -f oaitool
