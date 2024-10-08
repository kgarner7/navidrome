# Generating protobuf files

## Requirements

- protoc
- [go-plugin](https://github.com/knqyf263/go-plugin/releases)

## Building

```bash
protoc --go-plugin_out=./greeting --go-plugin_opt=paths=source_relative greeting.proto
```
