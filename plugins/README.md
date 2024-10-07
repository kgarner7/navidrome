Inspired by https://github.com/hashicorp/go-plugin/tree/main/examples/grpc

If you want to test with the plugin, pass in `KV_PLUGIN=...` to Navidrome before running.

Examples:

```bash
# Python
KV_PLUGIN="plugins/plugin-python-grpc/keyval_client.py" make dev
# GRPC (note, the output is named kv-go-grpc)
KV_PLUGIN=plugins/plugin-go-grpc/kv-go-grpc make dev
```
