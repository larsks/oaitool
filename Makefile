PKG=github.com/larsks/oaitool

GOSRC =  \
	 api/main.go \
	 api/host.go \
	 api/pullsecret.go \
	 api/datatypes.go \
	 api/cluster.go \
	 cli/host.go \
	 cli/main.go \
	 cli/cluster.go \
	 version/main.go \
	 main.go

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo unknown)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

GOLDFLAGS = \
	    -X '$(PKG)/version.BuildVersion=$(VERSION)' \
	    -X '$(PKG)/version.BuildRef=$(COMMIT)' \
	    -X '$(PKG)/version.BuildDate=$(DATE)'

all: oaitool

oaitool: $(GOSRC)
	go build -o $@ -ldflags "$(GOLDFLAGS)"

clean:
	rm -f oaitool
