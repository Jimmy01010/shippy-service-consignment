build:
	protoc --proto_path=. --go_out=. --micro_out=./ \
    		proto/consignment/consignment.proto
#	protoc --go_out=./proto/consignment \
#	  ./proto/consignment/consignment.proto
#	protoc --go-grpc_out=./proto/consignment \
#	  ./proto/consignment/consignment.proto

