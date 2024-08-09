gen-deps:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
	@go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

gen-grpc:
	@protoc -I=. -I=protos/ --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			--validate_out="lang=go,paths=source_relative:." \
			protos/auth/v1/auth.proto

gen: gen-deps gen-grpc
.PHONY: gen
