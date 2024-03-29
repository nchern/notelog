VERSION_FILE="pkg/cli/version_generated.go"

.PHONY: build
build: vet
	go build ./...

.PHONY: install
install: test gen-version
	go install ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test: build
	go test -race ./...

.PHONY: gen-version
gen-version:
	@./generate-version.sh $(VERSION_FILE)

.PHONY: coverage
coverage: vet
	@./coverage.sh

.PHONY: coverage-html
coverage-html: vet
	@./coverage.sh html
