.PHONY: out
out:
	go mod vendor
	go run .

.PHONY: docker
docker:
	docker build -t marnixkade .
