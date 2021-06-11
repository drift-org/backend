GO_BIN=~/go/bin

dev:
	${GO_BIN}/air
run:
	go build -o ./tmp/main .
	./tmp/main
clean: fmt tidy
fmt:
	go fmt ./...
tidy:
	go mod tidy

# Sets up our git hooks locally.
# Will remove any pre-existing hooks.
setup-hooks:
	find .git/hooks -type l -exec rm {} \;
	find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;
	chmod u+x .git/hooks/*

# Install packages referenced
install:
	go get -t

test:
	${GO_BIN}/ginkgo -p $(ARGS) ./...