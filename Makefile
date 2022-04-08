# init project path
HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/build
CONFDIR	:= $(HOMEDIR)/conf
DATADIR := $(HOMEDIR)/data

# init command params
GO      := ${GOROOT}/bin/go
GOROOT  := ${GOROOT}
GOPATH  := $(shell $(GO) env GOPATH)
GOMOD   := $(GO) mod
GOBUILD := $(GO) build -trimpath
GOTEST  := $(GO) test -gcflags="-N -l"
GOPKGS  := $$($(GO) list ./...| grep -vE "vendor")

# make, make all
all: prepare clean compile package

#make prepare, download dependencies
prepare: gomod

gomod:
	$(GOMOD) download

#make compile
compile: build

build:
	$(GOBUILD) -o $(HOMEDIR)/go-spider

# make test, test your code
test: prepare test-case
test-case:
	$(GOTEST) -v -cover $(GOPKGS)

# make package
package: package-bin package-conf package-data
package-bin:
	mkdir -p $(OUTDIR)/bin
	mv go-spider $(OUTDIR)/bin
package-conf:
	mkdir -p $(OUTDIR)/conf
	cp -r $(CONFDIR)/* $(OUTDIR)/conf
package-data:
	mkdir -p $(OUTDIR)/data
	cp -r $(DATADIR)/* $(OUTDIR)/data

# make clean
clean:
	rm -rf $(OUTDIR)
	rm -rf $(HOMEDIR)/go-spider

# avoid filename conflict and speed up build 
.PHONY: all prepare compile test package clean build
