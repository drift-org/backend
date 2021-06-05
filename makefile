# replace this variable with the location that leads to your air executable file
air_location=~/go/bin/air 
dev:
	${air_location}
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
	ginkgo -v ./...