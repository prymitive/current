GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN = $(shell go env GOPATH)/bin
endif

COVER_DIR     = .cover
COVER_PROFILE = $(COVER_DIR)/coverage.out

$(GOBIN)/golangci-lint: tools/golangci-lint/go.mod tools/golangci-lint/go.sum
	go install -modfile=tools/golangci-lint/go.mod github.com/golangci/golangci-lint/cmd/golangci-lint
.PHONY: lint
lint: $(GOBIN)/golangci-lint
	$(GOBIN)/golangci-lint run

$(GOBIN)/gofumpt: tools/gofumpt/go.mod tools/gofumpt/go.sum
	go install -modfile=tools/gofumpt/go.mod mvdan.cc/gofumpt
$(GOBIN)/goimports: tools/goimports/go.mod tools/goimports/go.sum
	go install -modfile=tools/goimports/go.mod golang.org/x/tools/cmd/goimports
.PHONY: format
format: $(GOBIN)/gofumpt $(GOBIN)/goimports
	$(GOBIN)/gofumpt -extra -l -w .
	$(GOBIN)/goimports -local github.com/cloudflare/pint -w .

.PHONY: test
test:
	mkdir -p $(COVER_DIR)
	echo 'mode: atomic' > $(COVER_PROFILE)
	go test \
		-covermode=atomic \
		-coverprofile=$(COVER_PROFILE) \
		-coverpkg=./... \
		-race \
		-count=5 \
		-timeout=15m \
		./...

.PHONY: benchmark
benchmark:
	go test -v -count=1 -run=none -bench=. -benchmem ./...
