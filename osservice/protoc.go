package osservice

//go:generate protoc --go_out=osproto --go_opt=paths=source_relative --go-grpc_out=osproto --go-grpc_opt=paths=source_relative --proto_path osproto time.proto store.proto net.proto wifi.proto
