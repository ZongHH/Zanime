import threading
from threading import Lock

class ReadWriteLock:
    """ 修复后的读写锁实现 """
    def __init__(self):
        self._lock = threading.Lock()  # 用于保护计数器
        self._readers_count = 0
        self._writer_lock = threading.Lock()  # 写锁
        self._readers_lock = threading.Lock()  # 读锁
        
    def acquire_read(self):
        with self._lock:
            self._readers_count += 1
            if self._readers_count == 1:
                self._writer_lock.acquire()  # 第一个读者需要获取写锁
    
    def release_read(self):
        with self._lock:
            self._readers_count -= 1
            if self._readers_count == 0:
                self._writer_lock.release()  # 最后一个读者释放写锁
    
    def acquire_write(self):
        self._readers_lock.acquire()  # 阻止新的读者
        self._writer_lock.acquire()   # 等待现有读者完成
    
    def release_write(self):
        self._writer_lock.release()   # 允许读者或其他写者
        self._readers_lock.release()  # 释放读者锁

class ReadLockContext:
    def __init__(self, lock):
        self.lock = lock
    
    def __enter__(self):
        self.lock.acquire_read()
        return self
        
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.lock.release_read()

class WriteLockContext:
    def __init__(self, lock):
        self.lock = lock
    
    def __enter__(self):
        self.lock.acquire_write()
        return self
        
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.lock.release_write()