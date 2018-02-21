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
	go get github.com/jstemmer/go-junit-report
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

index_subjects:
	mkdir -p $(DATA_DIR)
	cat $(CACHE)/Subjects.2014.part01.xml | $(WORKDIR)/a_indexer_$(OS)_$(ARCH) --dbpath=$(DATA_DIR)/marcdex.db

covert_fast_xml_to_binary:
	$(WORKDIR)/a_marc2marc_$(OS)_$(ARCH) -i $(CACHE)/FASTChronological.marcxml -f m > $(CACHE)/FASTChronological.marc
	$(WORKDIR)/a_marc2marc_$(OS)_$(ARCH) -i $(CACHE)/FASTCorporate.marcxml -f m > $(CACHE)/FASTCorporate.marc
	$(WORKDIR)/a_marc2marc_$(OS)_$(ARCH) -i $(CACHE)/FASTEvent.marcxml -f m > $(CACHE)/FASTEvent.marc
	$(WORKDIR)/a_marc2marc_$(OS)_$(ARCH) -i $(CACHE)/FASTFormGenre.marcxml -f m > $(CACHE)/FASTFormGenre.marc
	$(WORKDIR)/a_marc2marc_$(OS)_$(ARCH) -i $(CACHE)/FASTGeographic.marcxml -f m > $(CACHE)/FASTGeographic.marc
	$(WORKDIR)/a_marc2marc_$(OS)_$(ARCH) -i $(CACHE)/FASTPersonal.marcxml -f m > $(CACHE)/FASTPersonal.marc
	$(WORKDIR)/a_marc2marc_$(OS)_$(ARCH) -i $(CACHE)/FASTTitle.marcxml -f m > $(CACHE)/FASTTitle.marc
	$(WORKDIR)/a_marc2marc_$(OS)_$(ARCH) -i $(CACHE)/FASTTopical.marcxml -f m > $(CACHE)/FASTTopical.marc

index_fast:
	mkdir -p $(DATA_DIR)
	# $(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTChronological.marc --dbpath=$(DATA_DIR)/marcdex.db
	# $(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTCorporate.marc --dbpath=$(DATA_DIR)/marcdex.db
	# $(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTEvent.marc --dbpath=$(DATA_DIR)/marcdex.db
	$(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTFormGenre.marc --dbpath=$(DATA_DIR)/marcdex.db
	$(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTGeographic.marc --dbpath=$(DATA_DIR)/marcdex.db
	# $(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTPersonal.marc --dbpath=$(DATA_DIR)/marcdex.db
	# $(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTTitle.marc --dbpath=$(DATA_DIR)/marcdex.db
	$(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTTopical.marc --dbpath=$(DATA_DIR)/marcdex.db


time_index_fast:
	rm -fR $(DATA_DIR)
	mkdir -p $(DATA_DIR)
	time $(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/FASTFormGenre.marc --dbpath=$(DATA_DIR)/marcdex.db -cpuprofile=./profile.out


time_index_lcsh:
	rm -fR $(DATA_DIR)
	mkdir -p $(DATA_DIR)
	time $(WORKDIR)/a_indexer_$(OS)_$(ARCH) --inputpath=$(CACHE)/Subjects.2014.utf8.marc --dbpath=$(DATA_DIR)/marcdex.db # -cpuprofile=./profile.out
