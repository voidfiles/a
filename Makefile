PROJECT         :=a
CW              :=$(shell pwd)
GOFILES         :=$(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES      :=$(shell go list ./... | grep -v /vendor/| grep -v /checkers)
OS              := $(shell go env GOOS)
ARCH            := $(shell go env GOARCH)
CACHE           :=download-cache

BIN             := $(CW)/bin

GITHASH         :=$(shell git rev-parse --short HEAD)
GITBRANCH       :=$(shell git rev-parse --abbrev-ref HEAD)
BUILDDATE      	:=$(shell date -u +%Y%m%d%H%M)
GO_LDFLAGS		  ?= -s -w
GO_BUILD_FLAGS  :=-ldflags "${GOLDFLAGS} -X main.BuildVersion=${GITHASH} -X main.GitHash=${GITHASH} -X main.GitBranch=${GITBRANCH} -X main.BuildDate=${BUILDDATE}"
ARTIFACT_NAME   :=$(PROJECT)-$(GITHASH).tar.gz
ARTIFACT_DIR    :=$(PROJECT_DIR)/_artifacts
WORKDIR         :=$(PROJECT_DIR)/_workdir
DATA_DIR        :=$(CW)/data
MISC_DIR        :=$(CW)/_misc
TRIPLES_DIR     :=$(CW)/_triples
WORKDIR 	      :=$(CW)/_work
DATA_URL        :=http://id.loc.gov/static/data/authoritieschildrensSubjects.nt.skos.zip

# Determine commands by looking into cmd/*
COMMANDS=$(wildcard ${CW}/cmd/*)

# Determine binary names by stripping out the dir names
BINS=$(foreach cmd,${COMMANDS},$(notdir ${cmd}))

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(WORKDIR)/$(PROJECT)_linux_amd64 $(GO_BUILD_FLAGS)
	$(foreach BIN, $(BINS), (cd cmd/$(BIN) && CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(BIN)_linux_amd64 $(GO_BUILD_FLAGS));)

build:
	echo $(BINS)
	CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(OS)_$(ARCH) $(GO_BUILD_FLAGS)
	$(foreach BIN, $(BINS), (cd cmd/$(BIN) && CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(BIN)_$(OS)_$(ARCH) $(GO_BUILD_FLAGS));)


dependencies:
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/alecthomas/gometalinter
	go get github.com/golang/dep/cmd/dep
	go get github.com/stretchr/testify
	dep ensure
	gometalinter --install

lint:
	echo "metalinter..."
	gometalinter --enable=goimports --enable=unparam --enable=unused --disable=golint --disable=govet $(GOPACKAGES)
	echo "megacheck..."
	megacheck $(GOPACKAGES)
	echo "golint..."
	golint -set_exit_status $(GOPACKAGES)
	echo "go vet..."
	go vet --all $(GOPACKAGES)

init: dependencies

clean:
	rm -fR $(DATA_DIR)
	rm -fR $(BIN)
	rm -fR $(CACHE)
	rm -fR $(TRIPLES_DIR)
	rm -fR $(WORKDIR)

test:
	CGO_ENABLED=0 go test $(GOPACKAGES)

test-race:
	CGO_ENABLED=1 go test -race $(GOPACKAGES)

run_boltdb:
	$(WORKDIR)/a_$(OS)_$(ARCH) --db=bolt --dbpath=$(DATA_DIR)/cayley.db --indexpath=$(DATA_DIR)/search.db

run_indexer_boltdb:
	$(WORKDIR)/a_indexer_$(OS)_$(ARCH) --db=bolt --dbpath=$(DATA_DIR)/cayley.db  --indexpath=$(DATA_DIR)/search.db
