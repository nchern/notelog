.PHONY: build
build: vet
	go build ./...

.PHONY: install
install: test
	go get ./...

.PHONY: lint
lint:
	@golint ./...

.PHONY: vet
vet:
	 go vet ./...

.PHONY: test
test: build
	go test -race ./...
