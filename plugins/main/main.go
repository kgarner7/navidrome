//go:build tinygo.wasm

package main

import (
	"context"

	"github.com/navidrome/navidrome/plugins/greeting"
)

// main is required for TinyGo to compile to Wasm.
func main() {
	greeting.RegisterGreeter(MyPlugin{})
}

type MyPlugin struct{}

func (m MyPlugin) SayHello(ctx context.Context, request *greeting.GreetRequest) (*greeting.GreetReply, error) {
	return &greeting.GreetReply{
		Message: "Hello there, " + request.GetName(),
	}, nil

}
