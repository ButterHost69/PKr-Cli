protofiles:
	protoc ./proto/*.proto --go_out=. --go-grpc_out=.


protoc_files:
	protoc --go_out=. ./proto/*.proto
# Prefered One
# .PHONY protofiles