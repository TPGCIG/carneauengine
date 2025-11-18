BINARY_NAME=server.exe

.PHONY: build run dev clean

# Build the Go binary
build:
	go build -o $(BINARY_NAME) .

# Run the built binary
run:
	./$(BINARY_NAME)

# Build + run in one command (your dev workflow)
dev: build run

# Delete the binary
clean:
	del $(BINARY_NAME) 2>nul || true
