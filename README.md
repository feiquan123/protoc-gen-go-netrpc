# protoc-gen-go-netrpc

protoc's go-grpc plugin

### Last version install
```sh
go get github.com/feiquan123/protoc-gen-go-netrpc
go install $GOPATH/src/github.com/feiquan123/protoc-gen-go-netrpc
```

## Prev version install
```sh
go get github.com/feiquan123/protoc-gen-go-netrpc
go install -tags=prev $GOPATH/src/github.com/feiquan123/protoc-gen-go-netrpc
```

## Protoc gen go file

proto file `hello.proto`
```go
syntax = "proto3";

package api;

option go_package = "github.com/feiquan123/protoc-gen-go-netrpc/example/api";

message String {
    string value = 1;
}


service HelloService {
	rpc Hello(String) returns (String);
}
```

generate `hello.pb.go`
```sh
protoc --go-netrpc_out=plugins=netrpc:$GOPATH/src hello.proto
```

## [Example](https://github.com/feiquan123/protoc-gen-go-netrpc/blob/main/example/api/hello.proto)