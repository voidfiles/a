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
CAYLEY_URL      :=https://github.com/cayleygraph/cayley/releases/download/v0.7.0/cayley_0.7.0_$(OS)_$(ARCH).tar.gz
GO_LDFLAGS		  ?= -s -w
GO_BUILD_FLAGS  :=-ldflags "${GOLDFLAGS} -X main.BuildVersion=${GITHASH} -X main.GitHash=${GITHASH} -X main.GitBranch=${GITBRANCH} -X main.BuildDate=${BUILDDATE}"
CGO_LDFLAGS     :=-L/usr/local/opt/icu4c/lib
CGO_CFLAGS      :=-I/usr/local/opt/icu4c/include
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

CAYLEY_VERSION           := 0.7.0
CAYLEY_ARTIFACT_NAME     := cayley_$(CAYLEY_VERSION)_$(OS)_$(ARCH)
CAYLEY_FILENAME          := $(CAYLEY_ARTIFACT_NAME).tar.gz
CAYLEY_DOWNLOAD_URL      := https://github.com/cayleygraph/cayley/releases/download/v$(CAYLEY_VERSION)/$(CAYLEY_FILENAME)
CAYLEY_CMD               := $(BIN)/cayley

download_lc_authority_data:
	mkdir -p $(CACHE)
	$(MISC_DIR)/download_urls.sh $(MISC_DIR)/lc_authority_urls.txt $(CACHE)

download_small_data:
	mkdir -p $(CACHE)/
	cd $(CACHE) && curl -O -L $(DATA_URL)
	mkdir -p $(TRIPLES_DIR)/
	unzip "$(CACHE)/authorities*.zip" -d $(TRIPLES_DIR)/

unzip_lc_authority_data:
	mkdir -p $(TRIPLES_DIR)
	unzip "$(CACHE)/authorities*.zip" -d $(TRIPLES_DIR)/
	unzip "$(CACHE)/vocab*.zip" -d $(TRIPLES_DIR)/

load_triples:
	mkdir -p $(DATA_DIR)/cayley.db
	$(CAYLEY_CMD) init \
		--db bolt \
		--dbpath $(DATA_DIR)/cayley.db \
		-c $(CW)/conf/cayley.yml || true

	$(MISC_DIR)/load_triples.sh $(TRIPLES_DIR) $(CAYLEY_CMD) $(CW) $(DATA_DIR)/cayley.db

small_build: download_small_data load_triples

build-linux:
	CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(WORKDIR)/$(PROJECT)_linux_amd64 $(GO_BUILD_FLAGS)
	$(foreach BIN, $(BINS), (cd cmd/$(BIN) && CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(BIN)_linux_amd64 $(GO_BUILD_FLAGS));)

build:
	echo $(BINS)
	CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(OS)_$(ARCH) $(GO_BUILD_FLAGS)
	$(foreach BIN, $(BINS), (cd cmd/$(BIN) && CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(BIN)_$(OS)_$(ARCH) $(GO_BUILD_FLAGS));)


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

# install installs the cayley binary locally
install:
	echo "Installing cayley in $(CW)/bin/"
	## Create download cache if it doesn't exist...
	mkdir -p $(CW)/$(CACHE)
	## Fetch terraform and sha sums...
	(cd $(CW)/$(CACHE) && curl -O -L $(CAYLEY_DOWNLOAD_URL))
	## Make bin directory if it doesn't exist.
	mkdir -p $(BIN)
	## Unpack into the bin dir.
	(cd $(CW)/$(CACHE) && tar -xzf $(CW)/$(CACHE)/$(CAYLEY_FILENAME))
	mv $(CW)/$(CACHE)/$(CAYLEY_ARTIFACT_NAME)/cayley $(CW)/bin/
	echo -n "Installed cayley: " && $(BIN)/cayley version
	echo "Done..."

init: install

data: download_data unzip_data load_data

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

run_cayley:
	$(CAYLEY_CMD) http --host 127.0.0.1:64210 -d bolt -a $(DATA_DIR)/cayley.db

run_boltdb:
	$(WORKDIR)/a_$(OS)_$(ARCH) --db=bolt --dbpath=$(DATA_DIR)/cayley.db --indexpath=$(DATA_DIR)/search.db

run_indexer_boltdb:
	$(WORKDIR)/a_indexer_$(OS)_$(ARCH) --db=bolt --dbpath=$(DATA_DIR)/cayley.db  --indexpath=$(DATA_DIR)/search.db
