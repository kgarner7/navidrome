package host

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/plugins/grpc"
	"github.com/spf13/viper"
)

var (
	pluginViper = viper.New()
)

type pluginConfig struct {
	Name        string
	Description string
	AllowedUrls []string
	Config      map[string]string
}

func init() {
	pluginViper.SetConfigName("plugin")
	pluginViper.SetDefault("name", "")
	pluginViper.SetDefault("description", "")
	pluginViper.SetDefault("allowedurls", []string{})
}

type Plugin struct {
	ClosablePlugin
	Name        string
	Description string
	Directory   string
	Config      map[string]string
}

type ClosablePlugin interface {
	Close(ctx context.Context) error
	grpc.Greeter
}

// Adapted from https://stackoverflow.com/a/25747925/744298
func validateUrls(potentialUrls []string) error {
	for _, potential := range potentialUrls {
		u, err := url.Parse(potential)
		if err != nil {
			return err
		}

		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("bad url %s: must have scheme and host", potential)
		}
	}

	return nil
}

func LoadPlugin(ctx context.Context, dir string) (*Plugin, error) {
	directory := filepath.Dir(dir)

	pluginViper.SetConfigFile(filepath.Join(dir, "plugin.toml"))
	err := pluginViper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := pluginConfig{}
	if err = pluginViper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err = validateUrls(config.AllowedUrls); err != nil {
		return nil, err
	}

	hc := &http.Client{
		Timeout: consts.DefaultHttpClientTimeOut,
	}

	hostFuncs := NewHostFuncs(hc, directory, config.AllowedUrls)

	pluginWrapper, err := grpc.NewGreeterPlugin(ctx)
	if err != nil {
		return nil, err
	}

	plugin, err := pluginWrapper.Load(ctx, filepath.Join(dir, "plugin.wasm"), hostFuncs)
	if err != nil {
		return nil, err
	}

	_, err = plugin.Configure(ctx, &grpc.ConfigRequest{
		Config: config.Config,
	})

	if err != nil {
		plugin.Close(ctx)
		return nil, err
	}

	return &Plugin{
		Name:           config.Name,
		Description:    config.Description,
		Directory:      directory,
		ClosablePlugin: plugin,
	}, nil
}
