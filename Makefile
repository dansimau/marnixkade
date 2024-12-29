.PHONY: out
out:
	go run .

.PHONY: docker
docker:
	docker build -t marnixkade .
