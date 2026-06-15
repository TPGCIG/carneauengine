BINARY_NAME=server.exe
TICKET_TYPE_ID ?= 15

.PHONY: build run dev clean k6

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

# Reset DB + Redis state then run k6 spike test.
# Override ticket type with: make k6 TICKET_TYPE_ID=42
k6:
	powershell -ExecutionPolicy Bypass -File k6test.ps1 -TicketTypeId $(TICKET_TYPE_ID)
