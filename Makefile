GO := go
VERSION ?= $(shell git rev-parse HEAD)

word-find = $(word $2,$(subst -, ,$1))

./bin:
	mkdir bin

scredenziato-%: ./bin
	CGO_ENABLED=1 \
	GOOS=$(call word-find,$*,1) \
	GOARCH=$(call word-find,$*,2) \
	CC=$${CC_$(call word-find,$*,1)_$(call word-find,$*,2):-$$(which cc)} \
	CXX=$${CXX_$(call word-find,$*,1)_$(call word-find,$*,2):-$$(which c++)} \
	$(GO) build -v -ldflags="-X main.Version=$(VERSION)" -o bin/$@ ./cmd/scredenziato/...

build:
	make scredenziato-linux-amd64
	make scredenziato-darwin-amd64
	make scredenziato-darwin-arm64
	make scredenziato-windows-amd64

build-docker-image:
	if [ "z$USE_FLEXO_CACHE" != "z" ]; then \
		BUILD_ARGS="--build-arg=USE_FLEXO_CACHE=1"; \
	fi; \
	docker build $${BUILD_ARGS} -t scredenziato_builder:latest .

build-with-docker:
	docker run --rm -e VERSION=$(VERSION) -v $${PWD}:/domus/scredenziato -ti scredenziato_builder:latest

.PHONY: build
