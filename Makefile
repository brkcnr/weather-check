run:
	@trap 'echo "Server stopped successfully."; exit' INT; \
	go run cmd/server/main.go