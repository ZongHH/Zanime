import grpc
import traceback
from concurrent import futures
from src.grpc.message_pb2_grpc import RecommendServiceServicer, add_RecommendServiceServicer_to_server
from src.grpc.message_pb2 import RecommendationResponse, RecommendationResponseList
from src.calculate.similarity import recommender  # 确保这个路径正确

class RecommendService(RecommendServiceServicer):
    def __init__(self):
        # 确保推荐器已经加载了模型状态
        try:
            recommender.load_state()
        except FileNotFoundError:
            print("Warning: Model state file not found, please ensure that model training has been performed")

    def ListRecommendations(self, request, context):
        try:
            print(f"Received recommendation request: user_id={request.user_id}, Client Address: {context.peer()}")
            # 设置响应超时
            context.set_code(grpc.StatusCode.OK)
            print(f"recommend for user_id: {request.user_id}")
            
            # 调用推荐器获取推荐结果
            recommendations = recommender.recommend(request.user_id, n_recommendations=10)

            for item in recommendations:
                print(item)
            
            # 构建响应
            response_list = RecommendationResponseList()
            for rec in recommendations:
                response = RecommendationResponse()
                response.video_id = rec['video_id']
                response.video_name = rec['video_name']
                response.cover_image_url = rec['cover_image_url']
                response_list.recommendations.append(response)
            
            return response_list
            
        except Exception as e:
            error_stack = traceback.format_exc()
            print(f"Recommendation error: {str(e)}\nStack trace:\n{error_stack}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f'推荐服务出错: {str(e)}')
            return RecommendationResponseList()

def serve(port=50051):
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10),
        options=[            
            ('grpc.max_send_message_length', 10 * 1024 * 1024),
            ('grpc.max_receive_message_length', 10 * 1024 * 1024),
            ('grpc.keepalive_time_ms', 10000),
            ('grpc.keepalive_timeout_ms', 5000),
        ]
    )
    add_RecommendServiceServicer_to_server(RecommendService(), server)
    server.add_insecure_port(f'0.0.0.0:{port}')  # 修改为监听所有接口
    server.start()
    print(f"Recommendation service has been started on port {port}")
    return server
