.PHONY: build clean dist test vet fmt

build: # Build executable for local target
	go build

clean: # Clean cache
	go clean -testcache -modcache -cache

dist: # Build executable for supported platforms
	for ARCH in amd64 arm64 ppc64le; do \
		for OS in linux darwin windows; do \
			if GOARCH=$${ARCH} GOOS=$${OS} go build -o dist/distronaut-$${OS}-$${ARCH}; then \
				echo "built $${OS}-$${ARCH}"; \
			fi \
		done \
	done;

test: # Run test
	go clean -testcache
	go run tests/server.go &
	curl --retry 5 --retry-connrefused 0.0.0.0:3000/ready
	go test ./...
	curl 0.0.0.0:3000/stop

vet: # Vet code
	go vet

fmt: # Format code
	go fmt
