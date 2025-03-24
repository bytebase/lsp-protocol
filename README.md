# lsp-protocol

Containing the LSP protocol stubs in Golang.

## Usage

```bash
go run ./internal/generate -o .
```

## How it's working

The `generate` tool will clone the [go/tools](https://github.com/golang/tools) repository, and run the [`gopls/internal/protocol/generate/main.go`](https://github.com/golang/tools/blob/master/gopls/internal/protocol/generate/main.go) file to generate the LSP protocol stubs. The `generate` tools will copy the
`tsjson.go` and `tsprotocol.go` files to the specified output directory.
