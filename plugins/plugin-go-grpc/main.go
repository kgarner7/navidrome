package main

import (
	"errors"

	"github.com/hashicorp/go-plugin"
	"github.com/navidrome/navidrome/plugins/shared"
)

type KV struct {
	data map[string]string
}

func kg() *KV {
	kv := KV{
		data: make(map[string]string),
	}
	return &kv
}

func (k *KV) Put(key string, value []byte) error {
	k.data[key] = string(value)
	return nil
}

func (k *KV) Get(key string) ([]byte, error) {
	value, ok := k.data[key]
	if !ok {
		return nil, errors.New("not found")
	}

	return []byte(value), nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv": &shared.KVGRPCPlugin{Impl: kg()},
		},

		GRPCServer: plugin.DefaultGRPCServer,
	})
}
