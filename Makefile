.DEFAULT: all
.PHONY: all clean realclean deps

HOST=quay.io
NAMESPACE=weaveworks
DOCKER?=docker

# NB because this outputs absolute file names, you have to be careful
# if you're testing out the Makefile with `-W` (pretend a file is
# new); use the full path to the pretend-new file, e.g.,
#  `make -W $PWD/registry/registry.go`
godeps=$(shell go list -f '{{join .Deps "\n"}}' $1 | grep -v /vendor/ | xargs go list -f '{{if not .Standard}}{{ $$dep := . }}{{range .GoFiles}}{{$$dep.Dir}}/{{.}} {{end}}{{end}}')

PROMETHEUS_SQL_EXPORTER_DEPS:=$(call godeps,.)

all: build/.prometheus_sql_exporter.done

clean:
	go clean
	rm -rf ./build

realclean: clean
	rm -rf ./cache

build/.%.done: docker/Dockerfile.%
	mkdir -p ./build/docker/$*
	cp $^ ./build/docker/$*/
	${DOCKER} build -t ${HOST}/${NAMESPACE}/$* -f build/docker/$*/Dockerfile.$* ./build/docker/$*
	${DOCKER} tag ${HOST}/${NAMESPACE}/$* ${HOST}/${NAMESPACE}/$*:test
	touch $@

build/.prometheus_sql_exporter.done: build/prometheus_sql_exporter

build/prometheus_sql_exporter: $(PROMETHEUS_SQL_EXPORTER_DEPS)
build/prometheus_sql_exporter: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $(LDFLAGS) -ldflags "-X main.version=$(shell ./docker/image-tag)" main.go

deps:
	if [ -z $(shell which glide) ] ; then curl https://glide.sh/get | sh ; fi
	glide i
