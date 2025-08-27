# services/auth-service/Makefile
.PHONY: proto clean help

help:
	@echo "Available commands:"
	@echo "  make proto  - Generate proto files"
	@echo "  make clean  - Clean generated files"

proto:
	protoc --go_out=proto --go_opt=paths=source_relative \
	       --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
	       --micro_out=proto --micro_opt=paths=source_relative \
	       auth.proto

clean:
	@echo "Cleaning generated proto files..."
	rm -f proto/*.pb.go
	@echo "✅ Cleaned!"

# Test if protoc is installed
check:
	@which protoc > /dev/null || (echo "❌ protoc not found" && exit 1)
	@which protoc-gen-go > /dev/null || (echo "❌ protoc-gen-go not found" && exit 1)
	@which protoc-gen-go-grpc > /dev/null || (echo "❌ protoc-gen-go-grpc not found" && exit 1)
	@echo "✅ All proto tools are installed"
	@protoc --version