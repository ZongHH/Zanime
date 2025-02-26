import pandas as pd
import pickle
import threading
import time
import os
from sklearn.preprocessing import MinMaxScaler
from sklearn.metrics.pairwise import cosine_similarity
from sklearn.feature_extraction.text import TfidfVectorizer
from datetime import datetime
from tqdm import tqdm
from functools import wraps
from pkg.readwriteLock import ReadLockContext, WriteLockContext, ReadWriteLock

def singleton(cls):
    """
    单例装饰器，确保类只有一个实例
    
    Args:
        cls: 需要装饰的类
        
    Returns:
        装饰后的类，保证只返回同一个实例
    """
    instances = {}
    
    @wraps(cls)
    def get_instance(*args, **kwargs):
        if cls not in instances:
            instances[cls] = cls(*args, **kwargs)
        return instances[cls]
    
    return get_instance

@singleton
class VideoRecommender:
    """
    视频推荐系统类
    
    实现了基于协同过滤的视频推荐功能,包括:
    - 用户-视频交互矩阵的维护
    - 视频相似度矩阵的计算
    - 个性化推荐
    - 冷启动推荐
    - 增量更新
    
    主要属性:
        user_video_matrix: 用户-视频交互矩阵
        video_similarity_matrix: 视频相似度矩阵  
        video_metadata: 视频元数据
        scaler: 用于标准化观看时长的缩放器
        name_vectorizer: 用于处理视频名称的TF-IDF向量化器
    """
    
    _instance_lock = threading.Lock()  # 用于单例模式的锁
    _is_initialized = False  # 初始化标志

    def __init__(self):
        """初始化推荐器的各项属性"""
        # 原始数据
        self.user_video_matrix = None  # 用户-视频交互矩阵
        self.video_similarity_matrix = None  # 视频相似度矩阵
        self.video_metadata = None  # 视频元数据
        
        # 数据副本(用于增量更新)
        self._temp_user_video_matrix = None  # 临时用户-视频矩阵
        self._temp_video_metadata = None  # 临时视频元数据
        
        # 其他成员变量
        self.scaler = MinMaxScaler()  # 用于标准化观看时长
        self.name_vectorizer = TfidfVectorizer(stop_words='english')  # 用于处理视频名称
        self._rw_lock = ReadWriteLock()  # 原数据读写锁
        self._rw_temp_lock = ReadWriteLock()  # 副本读写锁
        self._update_thread = None  # 更新线程
        self._update_interval = 600  # 更新间隔(秒)
        self._stop_update = threading.Event()  # 停止更新的事件标志
        self.last_update = None  # 最后更新时间

    def _start_update_thread(self):
        """
        启动定时更新线程
        
        该线程会定期将临时数据更新到主数据中,并重新计算相似度矩阵
        """
        def update_loop():
            while not self._stop_update.is_set():
                try:
                    # 等待下一次更新
                    self._stop_update.wait(self._update_interval)

                    # 计算新的相似度矩阵
                    user_video_matrix, video_metadata, video_similarity = self._compute_video_similarity_temp()
                        
                    # 更新数据
                    with WriteLockContext(self._rw_lock):
                        self.video_metadata = video_metadata
                        self.user_video_matrix = user_video_matrix
                        self.video_similarity_matrix = video_similarity
                        self.last_update = datetime.now()
                        self.save_state()
                        
                except Exception as e:
                    print(f"update_loop error: {e}")
                    time.sleep(60)  # 发生错误时等待1分钟再试

        self._update_thread = threading.Thread(target=update_loop, daemon=True)
        self._update_thread.start()

    def incremental_update(self, new_data):
        """
        增量更新推荐系统的数据
        
        Args:
            new_data: 新的数据,可以是CSV文件路径或DataFrame对象
            
        功能:
            - 更新临时数据集
            - 处理新用户的观看记录
            - 更新视频元数据
            - 启动更新线程(如果尚未启动)
        """
        if isinstance(new_data, str):
            new_data = pd.read_csv(new_data)
            
        with WriteLockContext(self._rw_temp_lock):
            # 首次更新时初始化临时数据和scaler
            if self._temp_video_metadata is None:
                self._temp_video_metadata = self.video_metadata.copy() if self.video_metadata is not None else None
                self._temp_user_video_matrix = self.user_video_matrix.copy() if self.user_video_matrix is not None else None
                if not hasattr(self.scaler, 'min_') or self.scaler.min_ is None:
                    # 如果scaler未拟合，使用第一批数据进行拟合
                    self.scaler.fit(new_data[['StayTime']])
            
            # 更新临时数据
            try:
                # 更新视频元数据
                new_metadata = new_data.drop_duplicates('VideoID')[['VideoID', 'VideoName', 'Genre', 'Area', 'CoverImageUrl']]
                if self._temp_video_metadata is None:
                    self._temp_video_metadata = new_metadata
                else:
                    self._temp_video_metadata = pd.concat([self._temp_video_metadata, new_metadata]).drop_duplicates('VideoID')
                
                # 更新用户-视频矩阵
                if not hasattr(self.scaler, 'min_') or self.scaler.min_ is None:
                    # 如果scaler还未拟合，先进行拟合
                    self.scaler.fit(new_data[['StayTime']])
                    new_data['normalized_stayTime'] = self.scaler.transform(new_data[['StayTime']])
                else:
                    # scaler已拟合，直接transform
                    new_data['normalized_stayTime'] = self.scaler.transform(new_data[['StayTime']])

                temp_matrix = new_data.pivot_table(
                    index="UserID", 
                    columns="VideoID", 
                    values="normalized_stayTime", 
                    fill_value=0
                )
                
                if self._temp_user_video_matrix is None:
                    self._temp_user_video_matrix = temp_matrix
                else:
                    self._temp_user_video_matrix = pd.concat([
                        self._temp_user_video_matrix, 
                        temp_matrix
                    ]).groupby(level=0).sum()
                
            except Exception as e:
                print(f"incremental_update error: {e}")
                raise

        # 如果更新线程还没启动，启动它
        if self._update_thread is None or not self._update_thread.is_alive():
            self._start_update_thread()


    def _compute_video_similarity_temp(self):
        """
        使用临时数据计算视频相似度矩阵
        
        Returns:
            tuple: (user_video_matrix, video_metadata, video_similarity)
                - user_video_matrix: 更新后的用户-视频矩阵
                - video_metadata: 更新后的视频元数据
                - video_similarity: 计算得到的视频相似度矩阵
                
        计算过程:
            1. 基于用户行为计算基础相似度
            2. 计算视频名称的相似度
            3. 根据类型和地区信息调整相似度
        """
        # 获取读取锁，确保在计算时数据不会被修改
        with ReadLockContext(self._rw_temp_lock):
            # 使用copy()创建数据的副本，防止直接修改原始数据
            user_video_matrix = self._temp_user_video_matrix.copy()
            video_metadata = self._temp_video_metadata.copy()

        # 基础相似度（基于用户行为）
        base_similarity = cosine_similarity(user_video_matrix.T)
        video_similarity = pd.DataFrame(
            base_similarity,
            index=user_video_matrix.columns,
            columns=user_video_matrix.columns
        )

        # 计算名称相似度
        name_matrix = self.name_vectorizer.fit_transform(video_metadata['VideoName'])
        name_similarity = cosine_similarity(name_matrix)

        # 在视频元数据上进行相似度调整
        for idx1, video1 in tqdm(enumerate(video_metadata['VideoID']), 
                                desc="调整相似度", 
                                total=len(video_metadata['VideoID'])):
            for idx2, video2 in enumerate(video_metadata['VideoID']):
                genre1 = set(video_metadata.iloc[idx1]['Genre'].split(','))
                genre2 = set(video_metadata.iloc[idx2]['Genre'].split(','))
                genre_sim = len(genre1 & genre2) / len(genre1 | genre2) if genre1 | genre2 else 0

                area_sim = 1 if video_metadata.iloc[idx1]['Area'] == video_metadata.iloc[idx2]['Area'] else 0

                # 更新相似度矩阵
                video_similarity.loc[video1, video2] += (
                    0.15 * genre_sim +  # genre_weight
                    0.05 * area_sim +   # area_weight
                    0.15 * name_similarity[idx1, idx2]  # name_weight
                )

        return user_video_matrix, video_metadata, video_similarity


    def recommend(self, user_id, n_recommendations=5):
        """
        为指定用户生成推荐列表
        
        Args:
            user_id: 用户ID
            n_recommendations: 推荐数量,默认5个
            
        Returns:
            list: 推荐结果列表,每个元素包含视频信息
            
        处理流程:
            1. 检查模型是否初始化
            2. 对于新用户使用冷启动推荐
            3. 对于已有用户,基于其观看历史和视频相似度生成推荐
        """
        if not self._is_initialized:
            raise ValueError("模型尚未初始化")

        with ReadLockContext(self._rw_lock):
            if self.user_video_matrix is None or self.video_similarity_matrix is None:
                raise ValueError("需要先训练模型")
                
            if user_id not in self.user_video_matrix.index:
                return self._cold_start_recommend(n_recommendations)
            
            # 获取用户观看历史
            user_history = self.user_video_matrix.loc[user_id].copy()  # 创建副本避免数据竞争
            watched_videos = user_history[user_history > 0].index
            
            # 计算推荐分数
            scores = pd.Series(0.0, index=self.video_similarity_matrix.columns)
            for video_id in watched_videos:
                scores += self.video_similarity_matrix[video_id] * user_history[video_id]
                
            # 排除已观看的视频
            scores = scores.drop(watched_videos)
            
            # 返回推荐结果
            recommendations = scores.nlargest(n_recommendations)
            return self._format_recommendations(recommendations)

    def _cold_start_recommend(self, n_recommendations):
        """
        冷启动推荐方法
        
        Args:
            n_recommendations: 推荐数量
            
        Returns:
            list: 推荐结果列表
            
        说明:
            - 用于处理新用户的推荐
            - 基于视频的总体观看量进行推荐
        """
        with ReadLockContext(self._rw_lock):
            popular_videos = self.user_video_matrix.sum().nlargest(n_recommendations)
            return self._format_recommendations(popular_videos)

    def _format_recommendations(self, recommendations):
        """
        格式化推荐结果
        
        Args:
            recommendations: 推荐的视频ID和分数
            
        Returns:
            list: 包含完整视频信息的推荐列表
            
        说明:
            - 将视频ID转换为完整的视频信息
            - 包含视频名称、类型、地区等信息
        """
        result = []
        with ReadLockContext(self._rw_lock):
            for video_id, score in recommendations.items():
                video_info = self.video_metadata[self.video_metadata['VideoID'] == video_id].iloc[0]
                result.append({
                    'video_id': video_id,
                    'video_name': video_info['VideoName'],
                    'genre': video_info['Genre'],
                    'area': video_info['Area'],
                    'cover_image_url': video_info['CoverImageUrl'],
                    'score': score
                })
        return result

    def save_state(self, path='model/model_state.pkl'):
        """
        保存模型状态到文件
        
        Args:
            path: 保存路径,默认为'model/model_state.pkl'
            
        说明:
            保存的状态包括:
            - 用户-视频矩阵
            - 数据缩放器
            - 最后更新时间
            - 视频元数据
            - 视频相似度矩阵
        """
        state = {
            'matrix': self.user_video_matrix,
            'scaler': self.scaler,
            'last_update': self.last_update,
            'video_metadata': self.video_metadata,
            'video_similarity_matrix': self.video_similarity_matrix
        }
        with open(path, 'wb') as f:
            pickle.dump(state, f)
    
    def load_state(self, path='model/model_state.pkl'):
        """
        加载模型状态
        
        Args:
            path: 模型状态文件路径
            
        功能:
            - 线程安全的状态加载
            - 如果文件不存在则创建新的状态文件
            - 初始化所有必要的数据结构
        """
        if not self._is_initialized:
            with self._instance_lock:
                if not self._is_initialized:
                    try:
                        # 确保目录存在
                        os.makedirs(os.path.dirname(path), exist_ok=True)
                        
                        try:
                            with open(path, 'rb') as f:
                                state = pickle.load(f)
                            with WriteLockContext(self._rw_lock):
                                self.user_video_matrix = state['matrix']
                                # 确保scaler已经拟合过
                                if state['scaler'] and hasattr(state['scaler'], 'min_'):
                                    self.scaler = state['scaler']
                                else:
                                    self.scaler = MinMaxScaler()  # 创建新的scaler
                                self.last_update = state['last_update']
                                self.video_metadata = state['video_metadata']
                                self.video_similarity_matrix = state['video_similarity_matrix']
                        except FileNotFoundError:
                            print(f"Model state file not found, creating new file: {path}")
                            # 创建空的初始状态，包括未拟合的scaler
                            state = {
                                'matrix': None,
                                'scaler': MinMaxScaler(),  # 未拟合的scaler
                                'last_update': None,
                                'video_metadata': None,
                                'video_similarity_matrix': None
                            }
                            # 保存初始状态
                            with open(path, 'wb') as f:
                                pickle.dump(state, f)
                            # 设置初始状态
                            with WriteLockContext(self._rw_lock):
                                self.user_video_matrix = state['matrix']
                                self.scaler = state['scaler']
                                self.last_update = state['last_update']
                                self.video_metadata = state['video_metadata']
                                self.video_similarity_matrix = state['video_similarity_matrix']
                            
                        self._is_initialized = True
                        print("Model state initialization completed")
                        
                    except Exception as e:
                        print(f"Error during model state loading/creation: {str(e)}")
                        raise

    def stop(self):
        """
        停止推荐器
        
        功能:
            - 设置停止标志
            - 等待更新线程结束
        """
        self._stop_update.set()
        if self._update_thread and self._update_thread.is_alive():
            self._update_thread.join(timeout=5)
            
    def __del__(self):
        """
        析构函数
        
        功能:
            - 确保在对象销毁时停止所有线程
        """
        self.stop()

    def recommend_similar(self, video_name: str, genre: str, n_recommendations: int = 10) -> list:
        """
        基于视频名称和类型推荐相似的视频
        
        Args:
            video_name: 视频名称
            genre: 视频类型
            n_recommendations: 推荐数量，默认10个
            
        Returns:
            list: 推荐结果列表，每个元素包含视频信息
            
        功能:
            - 基于视频名称的文本相似度
            - 基于视频类型的相似度
            - 综合计算得到最终推荐结果
        """
        if not self._is_initialized:
            raise ValueError("模型尚未初始化")

        with ReadLockContext(self._rw_lock):
            if self.video_metadata is None:
                raise ValueError("视频元数据不存在")

            try:
                # 将输入的genre字符串转换为集合
                input_genre_set = set(genre.split(','))
                
                # 计算与所有视频的相似度
                similarities = []
                
                # 获取输入名称的TF-IDF向量
                input_name_vector = self.name_vectorizer.transform([video_name])
                
                # 获取所有视频名称的TF-IDF向量
                all_names_vector = self.name_vectorizer.transform(self.video_metadata['VideoName'])
                
                # 计算名称相似度
                name_similarities = cosine_similarity(input_name_vector, all_names_vector)[0]
                
                # 计算每个视频的综合相似度
                for idx, row in self.video_metadata.iterrows():
                    # 计算类型相似度
                    video_genre_set = set(row['Genre'].split(','))
                    genre_sim = len(input_genre_set & video_genre_set) / len(input_genre_set | video_genre_set) if input_genre_set | video_genre_set else 0
                    
                    # 综合相似度计算（可以调整权重）
                    total_similarity = (
                        0.6 * name_similarities[idx] +  # 名称相似度权重
                        0.4 * genre_sim                 # 类型相似度权重
                    )
                    
                    similarities.append({
                        'video_id': row['VideoID'],
                        'video_name': row['VideoName'],
                        'genre': row['Genre'],
                        'area': row['Area'],
                        'cover_image_url': row['CoverImageUrl'],
                        'score': total_similarity
                    })
                
                # 按相似度排序并返回前n个推荐
                recommendations = sorted(similarities, key=lambda x: x['score'], reverse=True)[:n_recommendations]
                
                return recommendations
                
            except Exception as e:
                print(f"相似推荐计算错误: {e}")
                raise

# 创建全局单例实例
recommender = VideoRecommender()

# 确保程序退出时停止更新线程
import atexit
atexit.register(recommender.stop)