GOLANG_VERSION := 1.19

test:
	docker run --rm -it -v "${PWD}:/app" -w /app golang:${GOLANG_VERSION} go test -v ./...
