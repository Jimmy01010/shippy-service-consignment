module consignment-cli

go 1.18

//require (
//	github.com/Jimmy01010/shippy-service-consignment/consignment-service/proto/consignment v0.0.0-20220622080602-19599edcac0b
//	google.golang.org/grpc v1.45.0
//)

replace github.com/Jimmy01010/shippy/shippy-service-consignment/proto/consignment => ../consignment-service/proto/consignment

require (
	github.com/Jimmy01010/shippy/shippy-service-consignment/proto/consignment v0.0.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	google.golang.org/grpc v1.47.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
