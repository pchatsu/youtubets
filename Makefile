NAME     := youtubets
VERSION  := v0.1.0
REVISION := $(shell git rev-parse --short HEAD)
MAINGO   := cmd/youtubets.go

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
NOVENDOR := $(shell go list ./... | grep -v vendor)

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME) cmd/youtubets/youtubets.go

.PHONY: install
install:
	go install $(LDFLAGS) cmd/youtubets/youtubets.go

.PHONY: cross-build
cross-build: deps
	set -e; \
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME) cmd/youtubets/youtubets.go; \
		done; \
	done

.PHONY: dist
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \; && \
	cd ..

.PHONY: release
release:
	git tag $(VERSION)
	git push origin $(VERSION)

.PHONY: test
test:
	go test -coverpkg=./... -v $(NOVENDOR)

.PHONY: ci-test
ci-test:
	go test -coverpkg=./... -coverprofile=coverage.txt -v ./...

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: deps
deps: dep
	dep ensure -v

.PHONY: update-deps
update-deps: dep
	dep ensure -update -v

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: build_test
build_test:
	go build $(LDFLAGS) -o bin/hoge cmd/youtubets/youtubets.go
