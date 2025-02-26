import nsq
import json
import time
import os
import pandas as pd
import concurrent.futures
from datetime import datetime
from src.calculate import recommender

class NSQClient:
    def __init__(self, topic, channel, nsqd_tcp_addresses=['127.0.0.1:5150']):
        """
        初始化NSQ客户端
        :param topic: NSQ主题
        :param channel: NSQ频道
        :param nsqd_tcp_addresses: NSQ服务器地址列表
        """
        self.topic = topic
        self.channel = channel
        self.nsqd_tcp_addresses = nsqd_tcp_addresses
        self.reader = None
        self.recommender = recommender
        self.is_running = False
        self.reconnect_interval = 5  # 重连间隔（秒）
        self.executor = concurrent.futures.ThreadPoolExecutor(max_workers=5)

        # 确保日志目录存在
        self.log_dir = 'logs'
        os.makedirs(self.log_dir, exist_ok=True)

        # 尝试加载推荐器状态
        try:
            self.recommender.load_state()
            print("nsq_client init successfully")
        except Exception as e:
            print(f"nsq_client init failed: {e}")

    def _log_message(self, message):
        """记录消息到日志文件"""
        timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        log_file = os.path.join(self.log_dir, 'nsq_messages.log')
        with open(log_file, 'a', encoding='utf-8') as f:
            f.write(f"[{timestamp}] {message}\n")

    def message_handler(self, message):
        self._distribute_message(message)
    
    def _distribute_message(self, message):
        """
        处理接收到的消息
        :param message: NSQ消息对象
        """
        try:
            # 解码消息内容
            message_data = json.loads(message.body.decode('utf-8'))
            self._log_message(f"收到消息: {message_data}")  # 添加日志记录

            # 根据消息类型处理
            message_type = message_data.get('msg_type', '')

            if message_type == 'user_behavior':
                self._handle_user_behavior(message_data)
            elif message_type == 'video_update':
                self._handle_video_update(message_data)
            elif message_type == 'recommendation_request':
                self._handle_recommendation_request(message_data)

            # 完成消息处理
            message.finish()            
        except json.JSONDecodeError:
            error_msg = "Message format error, unable to parse JSON"
            print(error_msg)
            self._log_message(error_msg)
            message.finish()
        except Exception as e:
            error_msg = f"An error occurred while processing the message: {e}"
            print(error_msg)
            self._log_message(error_msg)
            message.requeue()

    def _handle_user_behavior(self, data):
        """处理用户行为数据"""
        try:
            # 将用户行为数据转换为DataFrame格式
            behavior_data = pd.DataFrame([{
                'UserID': data['user_id'],
                'VideoID': data['video_id'],
                'VideoName': data['video_name'],
                'Genre': data['genre'],
                'Area': data['area'],
                'Release': data['release'],
                'StayTime': data['progress'],
                'CoverImageUrl': data['cover_image_url']
            }])

            # 增量更新推荐模型
            self.recommender.incremental_update(behavior_data)
            # print(f"_handle_user_behavior success user_id: {data['user_id']}")
            
        except Exception as e:
            print(f"_handle_user_behavior failed: {e}")
            raise

    def _handle_video_update(self, data):
        """处理视频更新数据"""
        try:
            # 更新视频元数据
            video_data = pd.DataFrame([{
                'VideoID': data['video_id'],
                'VideoName': data['video_name'],
                'Genre': data['genre'],
                'Area': data['area']
            }])
            
            # 更新推荐模型中的视频信息
            self.recommender.video_metadata = pd.concat([self.recommender.video_metadata, video_data]).drop_duplicates('VideoID')
            
            print(f"Video data updated successfully: {data['video_id']}")
            
        except Exception as e:
            print(f"Failed to process video update data: {e}")
            raise

    def _handle_recommendation_request(self, data):
        """处理推荐请求"""
        try:
            user_id = data['user_id']
            n_recommendations = data.get('n_recommendations', 5)
            
            # 获取推荐结果
            recommendations = self.recommender.recommend(
                user_id=user_id,
                n_recommendations=n_recommendations
            )
            
            print(f"为用户 {user_id} 生成推荐完成: {recommendations}")
            
            # TODO: 这里可以添加推荐结果的发送逻辑
            
        except Exception as e:
            print(f"Failed to process recommendation request: {e}")
            raise

    def connect(self):
        """建立NSQ连接"""
        try:
            # 使用nsq.Reader
            self.reader = nsq.Reader(
                topic=self.topic,
                channel=self.channel,
                message_handler=self.message_handler,
                nsqd_tcp_addresses=self.nsqd_tcp_addresses,
                max_in_flight=30,
                heartbeat_interval=30,
                timeout=5.0
            )
            self.is_running = True
            print(f"Connected to NSQ - Topic: {self.topic}, Channel: {self.channel}")
            nsq.run()
        except Exception as e:
            print(f"Connection to NSQ failed: {str(e)}")
            self.reconnect()

    def start(self):
        """启动NSQ客户端"""
        print("Start NSQ client")
        self.connect()

    def stop(self):
        """停止NSQ客户端"""
        if self.executor:
            self.executor.shutdown(wait=True)

        if self.reader:
            self.reader.close()
        self.is_running = False
        print("NSQ client has stopped")

    def reconnect(self):
        """重连机制"""
        if not self.is_running:
            print(f"Attempt to reconnect to NSQ, wait {self.reconnect_interval} seconds")
            time.sleep(self.reconnect_interval)
            try:
                if self.reader:
                    self.reader.close()
                self.connect()
            except Exception as e:
                print(f"Reconnection failed: {str(e)}")
                self.reconnect()

def test():
    """测试函数"""
    client = NSQClient(
        topic='user_video_behavior',
        channel='update',
        nsqd_tcp_addresses=['127.0.0.1:5150']
    )
    
    try:
        client.start()
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\n接收到退出信号")
        client.stop()
    except Exception as e:
        print(f"运行时发生错误: {str(e)}")
        client.stop()