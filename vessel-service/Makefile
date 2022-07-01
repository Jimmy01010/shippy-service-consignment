build:
	protoc --proto_path=. --go_out=proto/vessel/ \
	       --micro_out=proto/vessel \
    	proto/vessel/vessel.proto
	    # dep 工具暂不可用，直接手动编译
	    #GOOS=linux GOARCH=amd64 go build
	    #docker build -t vessel-service .
