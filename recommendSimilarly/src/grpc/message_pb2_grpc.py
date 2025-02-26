# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

import src.grpc.message_pb2 as message__pb2

GRPC_GENERATED_VERSION = '1.68.1'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in message_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class RecommendServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.ListRecommendations = channel.unary_unary(
                '/recommend.RecommendService/ListRecommendations',
                request_serializer=message__pb2.RecommendationRequest.SerializeToString,
                response_deserializer=message__pb2.RecommendationResponseList.FromString,
                _registered_method=True)


class RecommendServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def ListRecommendations(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_RecommendServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'ListRecommendations': grpc.unary_unary_rpc_method_handler(
                    servicer.ListRecommendations,
                    request_deserializer=message__pb2.RecommendationRequest.FromString,
                    response_serializer=message__pb2.RecommendationResponseList.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'recommend.RecommendService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('recommend.RecommendService', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class RecommendService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def ListRecommendations(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/recommend.RecommendService/ListRecommendations',
            message__pb2.RecommendationRequest.SerializeToString,
            message__pb2.RecommendationResponseList.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
