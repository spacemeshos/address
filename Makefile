generate:
	@echo "Generating wasm..."
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./build
	GOOS=js GOARCH=wasm go build -o build/address.wasm address.go

lint:
	golangci-lint run --new-from-rev=origin/master --config .golangci.yml

test:
	go test  ./...

test-race:
	go test -race  ./...