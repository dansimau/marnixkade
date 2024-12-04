.PHONY: out
out:
	go run main.go

.PHONY: docker
docker:
	docker build -t marnixkade .
