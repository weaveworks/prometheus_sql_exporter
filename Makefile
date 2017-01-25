.DEFAULT: all
.PHONY: all clean realclean deps integration build test

HOST=quay.io
NAMESPACE=weaveworks
DOCKER?=docker

# NB because this outputs absolute file names, you have to be careful
# if you're testing out the Makefile with `-W` (pretend a file is
# new); use the full path to the pretend-new file, e.g.,
#  `make -W $PWD/registry/registry.go`
godeps=$(shell go list -f '{{join .Deps "\n"}}' $1 | grep -v /vendor/ | xargs go list -f '{{if not .Standard}}{{ $$dep := . }}{{range .GoFiles}}{{$$dep.Dir}}/{{.}} {{end}}{{end}}')

PROSE_DEPS:=$(call godeps,.)

all: build/.prometheus_sql_exporter.done

clean:
	go clean
	rm -rf ./build

realclean: clean
	rm -rf ./cache

build/.%.done: docker/Dockerfile.%
	mkdir -p ./build/docker/$*
	cp $^ ./build/docker/$*/
	${DOCKER} build -t ${HOST}/${NAMESPACE}/$*:$(shell ./docker/image-tag) -f build/docker/$*/Dockerfile.$* ./build/docker/$*
	${DOCKER} images
	touch $@

build/.prometheus_sql_exporter.done: build/prose

build/prose: $(PROSE_DEPS)
build/prose: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $(LDFLAGS) -ldflags "-X main.version=$(shell ./docker/image-tag)" main.go

## Build lifecycle helpers

deps:
	go get github.com/Masterminds/glide
	glide i --force

build: build/prose

test:
	go test -v -race $(shell glide novendor)

integration: build/.mocks.done
	${DOCKER} run -d -p 15432:5432 --name integration-db integration-db
	until docker logs integration-db 2>1 | grep "PostgreSQL init process complete;" ; do sleep 1 ; done
	go test -v -race -tags integration -timeout 30s $(shell glide novendor) || { echo "Integration tests failed" >&2; if [ -z ${CI} ] ; then ${DOCKER} rm -f integration-db ; fi ; exit 1; }
	if [ -z ${CI} ] ; then ${DOCKER} rm -f integration-db ; fi ;

build/.mocks.done: ./mocks/Dockerfile.integration-db
	mkdir -p build
	${DOCKER} build -t integration-db -f $^ ./mocks
	touch $@

