build:
	go build -o bin/heim-cli main.go
install:
	mv ./bin/heim-cli ${GOBIN}
