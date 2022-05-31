package pb

// protoc protoc-gen-go protoc-gen-go-grpc in env PATH
//go:generate echo running protoc with $PWD $GOARCH $GOOS $GOPACKAGE $GOFILE $GOLINE $DOLLAR
//go:generate protoc $PWD/collect.proto --proto_path=$PWD --go_out=$PWD/
//go:generate protoc $PWD/collect.proto --proto_path=$PWD --go-grpc_out=$PWD/
