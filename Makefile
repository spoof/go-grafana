GO        ?= go
GOFMT     ?= $(GO)fmt
pkgs       = $$($(GO) list ./... | grep -v vendor)

.PHONY: all
all: build

.PHONY: build
build:
	$(GO) build github.com/spoof/go-grafana/client/...
	$(GO) build github.com/spoof/go-grafana/grafana/...
	$(GO) build github.com/spoof/go-grafana/pkg/...

.PHONY: verify
verify: checkformat vet

.PHONY: checkformat
checkformat:
	@echo ">> checking code format"
	! $(GOFMT) -d $$(find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

.PHONY: vet
vet:
	@echo ">> inspect source code"
	$(GO) vet  $(pkgs)

.PHONY: fmt
fmt:
	@echo ">> format code"
	$(GO) fmt $(pkgs)

.PHONY: test
test:
	@echo ">> running all tests"
	$(GO) test -v $(pkgs)