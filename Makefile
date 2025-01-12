build:
	go build -o bin/heimdahl main.go
install:
	mv ./bin/heimdahl ${GOBIN}
release:
	@echo ${GITHUB_TOKEN}
	@if [ -z "${GITHUB_TOKEN}" ]; then \
			echo "Error: GITHUB_TOKEN environment variable is not set"; \
			exit 1; \
	fi
	./scripts/build_and_release.sh ${GITHUB_TOKEN}
