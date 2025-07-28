GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN = $(shell go env GOPATH)/bin
endif

COVER_DIR     = .cover
COVER_PROFILE = $(COVER_DIR)/coverage.out

$(GOBIN)/golangci-lint: tools/golangci-lint/go.mod tools/golangci-lint/go.sum
	go install -modfile=tools/golangci-lint/go.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint
.PHONY: lint
lint: $(GOBIN)/golangci-lint
	$(GOBIN)/golangci-lint run
.PHONY: format
format: $(GOBIN)/golangci-lint
	$(GOBIN)/golangci-lint fmt

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
