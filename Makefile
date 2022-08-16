lint:
	golangci-lint run --new-from-rev=origin/master --config .golangci.yml

lint-fix:
	golangci-lint run --new-from-rev=master --config .golangci.yml --fix

test:
	go test  ./...

test-race:
	go test -race  ./...