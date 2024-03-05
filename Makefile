protofiles:
		cd client && protoc --proto_path=./myserver/proto ./myserver/proto/*.proto --go_out=. --go-grpc_out=.

.PHONY protofiles