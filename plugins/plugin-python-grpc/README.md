# Requirements

## Package requirements

```bash
# STRONGLY RECOMMENDED
virtualenv env
source env/bin/activate
pip install -r requirements.txt
```

## Build output files

```bash
python -m grpc_tools.protoc -I ../proto --python_out=. --pyi_out=. --grpc_python_out=. ../proto/example.proto
```

Notes:

- `--python_out` generates the `pb2.py` file
- `--pyi_out` generates the `pb2.pyi` file with typing
- `--grpc_python_out` generates the `pb2_grpc.py` file
