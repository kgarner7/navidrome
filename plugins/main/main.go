//go:build tinygo.wasm

package main

import (
	"context"
	"fmt"

	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/navidrome/navidrome/plugins/grpc"
)

// main is required for TinyGo to compile to Wasm.
func main() {
	grpc.RegisterGreeter(&MyPlugin{})
}

type MyPlugin struct {
	Config map[string]string
}

func (m *MyPlugin) SayHello(ctx context.Context, request *grpc.GreetRequest) (*grpc.GreetReply, error) {
	hostFunctions := grpc.NewHostFunctions()

	resp, _ := hostFunctions.Request(ctx, &grpc.HttpRequest{
		Url:    "https://example.org",
		Method: grpc.HttpMethod_GET,
	})

	hostFunctions.DoLog(ctx, &grpc.LogRequest{
		Log:   fmt.Sprint(m.Config),
		Level: grpc.LogLevel_Error,
	})

	if resp.Error != nil {
		hostFunctions.DoLog(ctx, &grpc.LogRequest{
			Log:   resp.Error.Message,
			Level: grpc.LogLevel_Error,
		})

		return &grpc.GreetReply{
			Message: "Fail" + resp.Error.Message,
		}, nil
	} else if resp.Response != nil {
		return &grpc.GreetReply{
			Message: "Hello there, " + request.GetName() + string(resp.Response.Response),
		}, nil
	} else {
		return &grpc.GreetReply{
			Message: "whut",
		}, nil
	}
}

func (m *MyPlugin) Configure(ctx context.Context, request *grpc.ConfigRequest) (*emptypb.Empty, error) {
	m.Config = request.Config
	return &emptypb.Empty{}, nil
}

var _ grpc.Greeter = (*MyPlugin)(nil)
