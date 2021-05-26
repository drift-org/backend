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