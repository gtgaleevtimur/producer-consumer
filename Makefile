.PHONY:

test:
	go clean -testcache && go test ./internal/consumer/core/ -v && go clean -testcache && go test ./internal/producer/core/ -v

run:
	docker compose up --build
