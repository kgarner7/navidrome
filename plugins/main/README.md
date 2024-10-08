# Go Plugin

## Requirements

- tinygo
- go1.19 - 1.22
- change the GO version/toolchain in go.mod to match the GO compiler you are using
- Have already built the protos

## Building

```bash
tinygo build -o main.wasm -scheduler=none -target=wasi --no-debug  main.go
```
