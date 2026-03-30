VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BINARY := bin/openspec

.PHONY: build test lint cover clean upstream-check-build upstream-check upstream-check-dry verify-build verify verify-audit verify-compat

build:
	@mkdir -p bin
	devbox run -- go build $(LDFLAGS) -o $(BINARY) ./cmd/openspec

test:
	devbox run -- go test ./... -count=1

lint:
	devbox run -- go vet ./...

cover:
	devbox run -- go test ./... -count=1 -coverprofile=coverage.out -covermode=atomic
	devbox run -- go tool cover -func=coverage.out

cover-html: cover
	devbox run -- go tool cover -html=coverage.out -o coverage.html

upstream-check-build:
	docker build -f containers/upstream-check/Containerfile -t openspec-upstream-check .

upstream-check: upstream-check-build
	docker run --rm \
		-e GH_TOKEN \
		-e GITHUB_REPOSITORY \
		-v "$(PWD)/UPSTREAM.md:/work/UPSTREAM.md:ro" \
		openspec-upstream-check \
		--repo "$${GITHUB_REPOSITORY}"

upstream-check-dry: upstream-check-build
	docker run --rm \
		-v "$(PWD)/UPSTREAM.md:/work/UPSTREAM.md:ro" \
		openspec-upstream-check \
		--dry-run

verify-build:
	docker build -f containers/verify/Containerfile -t openspec-verify .

verify: verify-build
	docker run --rm openspec-verify

verify-audit: verify-build
	docker run --rm openspec-verify audit

verify-compat: verify-build
	docker run --rm openspec-verify compat

clean:
	rm -f $(BINARY) coverage.out coverage.html
