import sys
import pickle
from src.server_manager import ServerManager
from src.config.config import config

server_manager = None

def clear_specific_data_in_pkl(path='model/model_state.pkl'):
    with open(path, 'rb') as f:
        state = pickle.load(f)

    # 清除部分字段数据
    state['matrix'] = None
    state['video_metadata'] = None
    state['video_similarity_matrix'] = None
    state['scaler'] = None
    state['last_update'] = None

    # 保存修改后的状态
    with open(path, 'wb') as f:
        pickle.dump(state, f)
    
    print(f"{path} 文件中的部分数据已被清空")


def signal_handler(sig, frame):
    if server_manager:
        server_manager.stop()
    sys.exit(0)

if __name__ == '__main__':
    try:
        print(f"gRPC配置: {config.get_grpc_config()}")
        print(f"NSQ配置: {config.get_nsq_config()}")

        # 创建并启动服务管理器
        # clear_specific_data_in_pkl()
        server_manager = ServerManager()
        
        # 启动所有服务
        server_manager.start_grpc(port=config.get_grpc_config()['port'])
        server_manager.start_nsq(topic=config.get_nsq_config()['topic'], channel=config.get_nsq_config()['channel'], nsqd_tcp_addresses=config.get_nsq_config()['nsqd_tcp_addresses'])

        if server_manager:
            server_manager.stop()
            
    except Exception as e:
        print(f"Service error occurred: {e}")
        if server_manager:
            server_manager.stop()
        sys.exit(1)