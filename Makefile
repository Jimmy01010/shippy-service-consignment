build:
	protoc --go_out=./proto/consignment \
	  ./proto/consignment/consignment.proto
	protoc --go-grpc_out=./proto/consignment \
	  ./proto/consignment/consignment.proto
