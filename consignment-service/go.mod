module consignment

go 1.18

replace github.com/Jimmy01010/shippy-service-consignment/consignment-service/proto/consignment => ./proto/consignment

require (
	github.com/Jimmy01010/shippy-service-consignment/consignment-service/proto/consignment v0.0.0-20220622080602-19599edcac0b
	golang.org/x/net v0.0.0-20210510120150-4163338589ed
	google.golang.org/grpc v1.47.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/sys v0.0.0-20210502180810-71e4cd670f79 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
