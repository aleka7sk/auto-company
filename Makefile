.PHONY: fmt test vet validate build check

fmt:
	gofmt -w .

test:
	go test ./...
	@node scripts/test-guard.mjs

vet:
	go vet ./...

validate:
	@node scripts/validate-plugin.mjs

build:
	go build -o bin/autoco ./cmd/autoco

check: fmt test vet validate build
