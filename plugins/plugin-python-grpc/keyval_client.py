#!/bin/env python3
from typing import Dict

from concurrent import futures
from sys import stdout
from time import sleep

from grpc import server
from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc


from example_pb2_grpc import *
from example_pb2 import *

state: Dict[str, str] = {}


class MyKVServicer(KVServicer):
    def Get(self, request: GetRequest, context):
        if request.key in state:
            result = GetResponse()
            result.value = state[request.key].encode()
            return result

        raise IndexError()

    def Put(self, request: PutRequest, context):
        state[request.key] = request.value.decode()
        return Empty()


def run():
    # We need to build a health service to work with go-plugin
    health = HealthServicer()
    health.set("plugin", health_pb2.HealthCheckResponse.ServingStatus.Value("SERVING"))

    grpc_server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_KVServicer_to_server(MyKVServicer(), grpc_server)
    health_pb2_grpc.add_HealthServicer_to_server(health, grpc_server)
    grpc_server.add_insecure_port("127.0.0.1:1234")
    grpc_server.start()

    # Output information
    print("1|1|tcp|127.0.0.1:1234|grpc")
    stdout.flush()

    try:
        while True:
            sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        grpc_server.stop(0)


if __name__ == "__main__":
    run()
