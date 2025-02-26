import os
import yaml
from typing import List, Dict, Any

class Config:
    def __init__(self):
        self.grpc_config: Dict[str, Any] = {}
        self.nsq_config: Dict[str, Any] = {}
        self.load_config()

    def load_config(self) -> None:
        """
        加载配置文件
        配置优先级: 环境变量 > 配置文件
        """
        # 获取配置文件路径
        config_path = os.getenv('CONFIG_PATH', 'configs/config.yaml')
        
        # 加载YAML配置
        try:
            with open(config_path, 'r', encoding='utf-8') as f:
                config = yaml.safe_load(f)
                
            # gRPC配置
            self.grpc_config = {
                'host': os.getenv('GRPC_HOST', config['grpc']['host']),
                'port': int(os.getenv('GRPC_PORT', config['grpc']['port'])),
                'max_workers': int(os.getenv('GRPC_MAX_WORKERS', config['grpc']['max_workers']))
            }
            
            # NSQ配置
            self.nsq_config = {
                'nsqd_tcp_addresses': os.getenv('NSQ_ADDRESSES', 
                    config['nsq']['consumer']['nsqd_tcp_addresses']).split(',') 
                    if isinstance(os.getenv('NSQ_ADDRESSES'), str) 
                    else config['nsq']['consumer']['nsqd_tcp_addresses'],
                'topic': os.getenv('NSQ_TOPIC', config['nsq']['consumer']['topic']),
                'channel': os.getenv('NSQ_CHANNEL', config['nsq']['consumer']['channel']),
                'max_in_flight': int(os.getenv('NSQ_MAX_IN_FLIGHT', 
                    config['nsq']['consumer']['max_in_flight'])),
                'heartbeat_interval': int(os.getenv('NSQ_HEARTBEAT_INTERVAL', 
                    config['nsq']['consumer']['heartbeat_interval'])),
                'message_timeout': int(os.getenv('NSQ_MESSAGE_TIMEOUT', 
                    config['nsq']['consumer']['message_timeout']))
            }
            
        except FileNotFoundError:
            print(f"配置文件 {config_path} 未找到，使用默认配置")
            self._set_default_config()
        except Exception as e:
            print(f"加载配置文件出错: {str(e)}，使用默认配置")
            self._set_default_config()

    def _set_default_config(self) -> None:
        """设置默认配置"""
        self.grpc_config = {
            'host': '0.0.0.0',
            'port': 50051,
            'max_workers': 10
        }
        
        self.nsq_config = {
            'nsqd_tcp_addresses': ['127.0.0.1:4150'],
            'topic': 'video_events',
            'channel': 'recommendation',
            'max_in_flight': 5,
            'heartbeat_interval': 30,
            'message_timeout': 60
        }

    def get_grpc_config(self) -> Dict[str, Any]:
        """获取gRPC配置"""
        return self.grpc_config

    def get_nsq_config(self) -> Dict[str, Any]:
        """获取NSQ配置"""
        return self.nsq_config

# 创建全局配置实例
config = Config()
