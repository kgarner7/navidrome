package host

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/plugins/grpc"
)

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type pluginHostFuncs struct {
	allowedUrls []string
	name        string
	hc          httpDoer
}

func NewHostFuncs(hc httpDoer, name string, urls []string) *pluginHostFuncs {
	return &pluginHostFuncs{
		allowedUrls: urls,
		name:        name,
		hc:          hc,
	}
}

var (
	levelToFunc = map[grpc.LogLevel]func(...interface{}){
		grpc.LogLevel_Error: log.Error,
		grpc.LogLevel_Warn:  log.Warn,
		grpc.LogLevel_Info:  log.Info,
		grpc.LogLevel_Debug: log.Debug,
		grpc.LogLevel_Trace: log.Trace,
	}
)

func (p *pluginHostFuncs) DoLog(ctx context.Context, logReq *grpc.LogRequest) (*emptypb.Empty, error) {
	if logReq.Level == grpc.LogLevel_Unknown {
		return &emptypb.Empty{}, nil
	}

	f := levelToFunc[logReq.Level]
	f(ctx, logReq.Log, "plugin", p.name)
	return &emptypb.Empty{}, nil
}

func httpError(message string, errorType grpc.ErrorType) (*grpc.RpcHttpResponse, error) {
	return &grpc.RpcHttpResponse{
		Error: &grpc.RpcError{
			Message: message,
			Type:    errorType,
		},
	}, nil
}

// Make HTTP request to allowed domains
func (p *pluginHostFuncs) Request(ctx context.Context, pluginReq *grpc.HttpRequest) (*grpc.RpcHttpResponse, error) {
	targetUrl := pluginReq.Url
	ok := false
	for _, url := range p.allowedUrls {
		if strings.HasPrefix(targetUrl, url) {
			ok = true
			break
		}
	}

	if !ok {
		return httpError(fmt.Sprintf("bad url %s", targetUrl), grpc.ErrorType_BadAccess)
	}

	log.Trace(ctx, "requesting url", "plugin", p.name, "url", pluginReq.Url)

	req, err := http.NewRequestWithContext(ctx, grpc.HttpMethod_name[int32(pluginReq.Method)], pluginReq.Url, strings.NewReader(pluginReq.GetBody()))

	if err != nil {
		return httpError(err.Error(), grpc.ErrorType_Other)
	}

	resp, err := p.hc.Do(req)
	if err != nil {
		return httpError(err.Error(), grpc.ErrorType_Other)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return httpError(err.Error(), grpc.ErrorType_Other)
	}

	return &grpc.RpcHttpResponse{
		Response: &grpc.HttpResponse{
			Code:     int32(resp.StatusCode),
			Response: data,
		},
	}, nil
}

var _ grpc.HostFunctions = (*pluginHostFuncs)(nil)
