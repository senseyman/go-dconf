export DOCKER_BUILDKIT=1

GOLINT := golangci-lint

dep:
	go mod tidy
	go mod download

dep-update:
	go get -t -u ./...

dep-all: dep-update dep

lint: dep
	$(GOLINT) run --timeout=5m -c .golangci.yml

test:
	go test -tags=unit,integration -cover -race -count=1 -timeout=60s ./...
