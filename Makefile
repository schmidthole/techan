files := $(shell find . -name "*.go" | grep -v vendor)

bootstrap:
	go install -v golang.org/x/lint/golint@latest
	go install -v golang.org/x/tools/...@latest
	go install -v honnef.co/go/tools/cmd/staticcheck@latest

imports:
	goimports -w $(files)

clean:
	rm -f cover.out

test:
	go test

lint:
	golint -set_exit_status
	golint -set_exit_status example
	staticcheck github.com/schmidthole/techan
	staticcheck github.com/schmidthole/techan/example

bench: clean
	go test -bench .

coverage: clean
	go test . -coverprofile=cover.out
	go tool cover -func=cover.out
