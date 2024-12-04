.PHONY: docker
docker:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/marnixkade main.go
	docker build -t marnixkade .
	rm -rf build
