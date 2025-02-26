import threading
import time
from .grpc import serve
from .consumer import NSQClient

class ServerManager:
    def __init__(self):
        self.grpc_server = None
        self.nsq_client = None
        self.grpc_thread = None
        self.nsq_thread = None
        self._stop_event = threading.Event()

    def start_grpc(self, port=50051):
        """在新线程中启动gRPC服务器"""
        def run_grpc():
            self.grpc_server = serve(port=port)
            while not self._stop_event.is_set():
                time.sleep(1)
            
        self.grpc_thread = threading.Thread(target=run_grpc, name="gRPC-Server")
        self.grpc_thread.start()

    def start_nsq(self, topic, channel, nsqd_tcp_addresses):
        """在主线程中启动NSQ客户端"""
        self.nsq_client = NSQClient(topic, channel, nsqd_tcp_addresses)
        self.nsq_client.start()

    def stop(self):
        """停止所有服务"""
        self._stop_event.set()
        
        if self.nsq_client:
            self.nsq_client.stop()
            if self.nsq_thread:
                self.nsq_thread.join(timeout=5)
        
        if self.grpc_server:
            print("stop gRPC server")
            self.grpc_server.stop(0)
            if self.grpc_thread:
                self.grpc_thread.join(timeout=5)
        
        print("all services have been stopped")
