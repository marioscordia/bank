run:
		@go run ./cmd/
build:
	docker build -t bank .
start:
	docker run -d -p 8080:8080 --name bank bank