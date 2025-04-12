- 安装依赖

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

- 编译 `proto` 文件
```shell
protoc --go_out=. --go-grpc_out=. proto/hello.proto
```

- 启动服务
```shell
go run main.go
```
