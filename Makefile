.PHONY: lint
lint:
	golangci-lint run -v --fix

.PHONY: build-image
build-image:
	docker build --rm -t allisson/rinha-de-backend-2024-q1-golang .

.PHONY: run-server
run-server:
	go run cmd/rinha/main.go
