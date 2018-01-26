PROJECT         :=a
CW              :=$(shell pwd)
GOFILES         :=$(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES      :=$(shell go list ./... | grep -v /vendor/| grep -v /checkers)
OS              := $(shell go env GOOS)
ARCH            := $(shell go env GOARCH)

GITHASH         :=$(shell git rev-parse --short HEAD)
GITBRANCH       :=$(shell git rev-parse --abbrev-ref HEAD)
BUILDDATE      	:=$(shell date -u +%Y%m%d%H%M)
GO_LDFLAGS		  ?= -s -w
GO_BUILD_FLAGS  :=-ldflags "${GOLDFLAGS} -X main.BuildVersion=${GITHASH} -X main.GitHash=${GITHASH} -X main.GitBranch=${GITBRANCH} -X main.BuildDate=${BUILDDATE}"
ARTIFACT_NAME   :=$(PROJECT)-$(GITHASH).tar.gz
ARTIFACT_DIR    :=$(PROJECT_DIR)/_artifacts
WORKDIR         :=$(PROJECT_DIR)/_workdir
DATA_DIR        :=$(CW)/data
DATA_URL        :=http://id.loc.gov/static/data/authoritieschildrensSubjects.nt.zip

download_data:
	mkdir -p $(DATA_DIR)
	wget $(DATA_URL) --output-document=$(DATA_DIR)/authoritieschildrensSubjects.nt.zip

unzip_data:
	cd $(DATA_DIR) && unzip authoritieschildrensSubjects.nt.zip

load_data:
	mkdir -p $(DATA_DIR)/cayley.db
	docker-compose run --rm \
		--entrypoint "cayley load" \
		cayley \
		--init \
		-c /conf/cayley.yml \
		-i /data/authoritieschildrensSubjects.nt \
		--alsologtostderr=true

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(WORKDIR)/$(PROJECT)_linux_amd64 $(GO_BUILD_FLAGS)

build:
	CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(OS)_$(ARCH) $(GO_BUILD_FLAGS)


dependencies:
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/alecthomas/gometalinter
	go get github.com/golang/dep/cmd/dep
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

init: download_data unzip_data load_data

run:
	docker-compose up -d

stop:
	docker-compose down

logs:
	docker-compose logs -f --tail=200

clean:
	rm -fR $(DATA_DIR)

test:
	CGO_ENABLED=0 go test $(GOPACKAGES)

test-race:
	CGO_ENABLED=1 go test -race $(GOPACKAGES)
