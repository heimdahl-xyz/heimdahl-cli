.PHONY: build


OUTPUT_NAME = heimdahl-${GOOS}-${GOARCH}
ifeq ($(GOOS),windows)
    OUTPUT_NAME := $(OUTPUT_NAME).exe
endif

clean:
	rm -rf build

build: clean
	go build -o bin/heimdahl main.go
install:
	mv ./bin/heimdahl ${GOBIN}
release:
	@if [ -z "${GITHUB_TOKEN}" ]; then \
			echo "Error: GITHUB_TOKEN environment variable is not set"; \
			exit 1; \
	fi
	ls -l scripts/
	scripts/release.sh

binary: clean
	CGO_ENABLED=0 GOGC=off GOOS=${GOOS} GOARCH=${GOARCH} \
	go build -installsuffix nocgo -o "./build/${OUTPUT_NAME}" main.go
