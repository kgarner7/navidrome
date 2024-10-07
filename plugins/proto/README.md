# Requirements

## Base requirements

You should already have protoc/protobuf installed on your system.

## Protoc plugins

```bash
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.34.2
```

## Build output files

```bash
cd plugins/proto
protoc -I . example.proto --go_out=. --go-grpc_out=
```

Notes:

- `--go_out` generates the `.pb.go`
- `--go-grpc_out` generates the `_grpc.pb.go` used for GRPC
