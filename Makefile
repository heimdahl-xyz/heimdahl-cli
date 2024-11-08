build:
	go build -o bin/heimdahl main.go
install:
	mv ./bin/heimdahl ${GOBIN}
