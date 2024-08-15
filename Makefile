files := $(shell find . -name "*.go" | grep -v vendor)

bootstrap:
	go install -v golang.org/x/lint/golint@latest
	go install -v golang.org/x/tools/...@latest
	go install -v honnef.co/go/tools/cmd/staticcheck@latest

clean:
	goimports -w $(files)
	rm cover.out

test: clean
	go test

lint:
	golint -set_exit_status
	golint -set_exit_status example
	staticcheck github.com/schmidthole/techan
	staticcheck github.com/schmidthole/techan/example

bench: clean
	go test -bench .

coverage:
	go test . -coverprofile=cover.out
	go tool cover -func=cover.out
