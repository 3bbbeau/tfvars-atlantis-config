set dotenv-load := true

tool-golangci:
    @hash golangci-lint > /dev/null 2>&1; if [ $? -ne 0 ]; then \
    GOBIN="$(pwd)/tools/bin" go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
    fi

lint: tool-golangci
    @$(pwd)/tools/bin/golangci-lint run -E gofumpt --timeout 1m

build:
    @go build -ldflags="-X 'github.com/3bbbeau/tfvars-atlantis-config/cmd.v=${VERSION}'"

test:
    @go test -cover ./...
    @go clean --testcache
    @go test ./... -v
