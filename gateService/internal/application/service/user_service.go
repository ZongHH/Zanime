package service

import (
	"context"
	"errors"
	"fmt"
	"gateService/internal/domain/entity"
	"gateService/internal/domain/repository"
	"gateService/internal/infrastructure/config"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/password"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type UserServiceImpl struct {
	storageConfig         *config.StorageConfig
	userRepository        repository.UserRepository
	postRepository        repository.PostRepository
	postCommentRepository repository.PostCommentRepository
	jwtManager            *auth.JWTManager
	cookieManager         *auth.CookieManager
}

func NewUserServiceImpl(
	storageConfig *config.StorageConfig,
	userRepository repository.UserRepository,
	postRepository repository.PostRepository,
	postCommentRepository repository.PostCommentRepository,
	jwtManager *auth.JWTManager,
	cookieManager *auth.CookieManager,
) *UserServiceImpl {
	return &UserServiceImpl{
		storageConfig:         storageConfig,
		userRepository:        userRepository,
		postRepository:        postRepository,
		postCommentRepository: postCommentRepository,
		jwtManager:            jwtManager,
		cookieManager:         cookieManager,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, user *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// 检查用户是否已存在
	exist, err := s.userRepository.IsExistUser(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("用户已存在")
	}

	// 检查密码长度
	if !password.CheckPasswordLength(user.Password) {
		return nil, errors.New("密码长度必须在8-16位之间")
	}

	// 检查密码复杂度
	if !password.CheckPasswordComplexity(user.Password) {
		return nil, errors.New("密码必须包含大小写字母、数字和特殊字符")
	}

	// 检查邮箱格式
	if !password.IsValidEmail(user.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 创建用户
	_, err = s.userRepository.CreateUser(ctx, &entity.UserInfo{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		Code:    200,
		Message: "用户注册成功",
	}, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, user *dto.LoginRequest) (*dto.LoginResponse, error) {
	userInfo := &entity.UserInfo{
		Email:    user.Email,
		Password: user.Password,
	}
	exist, err := s.userRepository.VerifyUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("账号或密码错误")
	}

	token, err := s.jwtManager.GenerateToken(userInfo)
	if err != nil {
		return nil, err
	}

	if c, ok := ctx.(*gin.Context); ok {
		s.cookieManager.SetTokenCookie(c, token)
	}

	return &dto.LoginResponse{
		Code:     200,
		UserInfo: userInfo,
	}, nil
}

func (s *UserServiceImpl) Logout(ctx context.Context) error {
	if c, ok := ctx.(*gin.Context); ok {
		s.cookieManager.ClearTokenCookie(c)
	}
	return nil
}

func (s *UserServiceImpl) VerifyUser(ctx context.Context, user *dto.VerifyUserRequest) (*dto.VerifyUserResponse, error) {
	if c, ok := ctx.(*gin.Context); ok {
		userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
		if userInfo.UserID != user.UserID {
			return nil, errors.New("用户ID不匹配")
		}
	} else {
		return nil, errors.New("context is not a gin context")
	}
	return &dto.VerifyUserResponse{
		Code:    200,
		Message: "用户ID匹配",
	}, nil
}

func (s *UserServiceImpl) GetUserInfo(ctx context.Context, request *dto.UserInfoRequest) (*dto.UserInfoResponse, error) {
	userInfo, err := s.userRepository.GetUserByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	return &dto.UserInfoResponse{
		Code: 200,
		User: dto.UserInfo{
			ID:        userInfo.UserID,
			Username:  userInfo.Username,
			AvatarURL: userInfo.AvatarURL,
		},
	}, nil
}

func (s *UserServiceImpl) GetUserProfile(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.GetUserProfileResponse, error) {
	userInfo, err := s.userRepository.GetUserByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 获取用户详细信息
	userProfile := &dto.GetUserProfileResponse{
		Code: 200,
		Profile: &dto.UserDetail{
			UserID:        userInfo.UserID,
			Username:      userInfo.Username,
			Email:         userInfo.Email,
			AvatarURL:     userInfo.AvatarURL,
			IsVIP:         true,
			VIPLevel:      10,
			VIPExpireDate: "2099-01-01",
			RegisterTime:  userInfo.CreatedAt,
			Signature:     userInfo.Signature,
			Gender:        userInfo.Gender,
		},
	}

	return userProfile, nil
}

func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, user *dto.UpdateUserProfileRequest) (*dto.UpdateUserProfileResponse, error) {
	userInfo := &entity.UserInfo{
		UserID:    user.UserID,
		Username:  user.Username,
		Gender:    user.Gender,
		Signature: user.Signature,
		AvatarURL: user.AvatarURL,
	}

	// 更新用户信息
	err := s.userRepository.UpdateUser(ctx, userInfo)
	if err != nil {
		return nil, fmt.Errorf("更新用户信息失败: %v", err)
	}

	return &dto.UpdateUserProfileResponse{
		Code: 200,
	}, nil
}

func (s *UserServiceImpl) GetUserStats(ctx context.Context, user *dto.UserStatsRequest) (*dto.UserStatsResponse, error) {
	var (
		postCount     int
		favoriteCount int
		commentCount  int
		postErr       error
		favoriteErr   error
		commentErr    error
	)

	// 使用WaitGroup等待所有goroutine完成
	var wg sync.WaitGroup
	wg.Add(3)

	// 并发获取用户发布的帖子数量
	go func() {
		defer wg.Done()
		postCount, postErr = s.postRepository.GetUserPostCount(ctx, user.UserID)
	}()

	// 并发获取用户收藏的帖子数量
	go func() {
		defer wg.Done()
		favoriteCount, favoriteErr = s.postRepository.GetUserFavoritePostCount(ctx, user.UserID)
	}()

	// 并发获取用户评论数量
	go func() {
		defer wg.Done()
		commentCount, commentErr = s.postCommentRepository.GetUserCommentCount(ctx, user.UserID)
	}()

	// 等待所有goroutine完成
	wg.Wait()

	// 检查是否有错误发生
	if postErr != nil {
		return nil, fmt.Errorf("获取用户帖子数量失败: %v", postErr)
	}
	if favoriteErr != nil {
		return nil, fmt.Errorf("获取用户收藏帖子数量失败: %v", favoriteErr)
	}
	if commentErr != nil {
		return nil, fmt.Errorf("获取用户评论数量失败: %v", commentErr)
	}

	response := &dto.UserStatsResponse{
		Code: 200,
	}
	response.Data.FollowingCount = 0
	response.Data.PostCount = postCount
	response.Data.FavoritePostCount = favoriteCount
	response.Data.CommentCount = commentCount

	return response, nil
}

func (s *UserServiceImpl) UploadAvatar(ctx context.Context, user *dto.UploadAvatarRequest) (*dto.UploadAvatarResponse, error) {
	// 检查文件是否存在
	if user.Avatar == nil {
		return nil, fmt.Errorf("未上传头像文件")
	}

	// 打开文件
	file, err := user.Avatar.Open()
	if err != nil {
		return nil, fmt.Errorf("打开头像文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取头像文件失败: %v", err)
	}

	// 检查文件大小
	if len(fileBytes) > s.storageConfig.Avatar.MaxSize {
		return nil, fmt.Errorf("文件大小超过限制,最大允许%d MB", s.storageConfig.Avatar.MaxSize/1024/1024)
	}

	// 检查文件类型
	fileType := http.DetectContentType(fileBytes)
	validType := false
	for _, allowedType := range s.storageConfig.Avatar.AllowedTypes {
		if fileType == allowedType {
			validType = true
			break
		}
	}
	if !validType {
		return nil, fmt.Errorf("不支持的文件类型: %s, 仅支持: %v", fileType, s.storageConfig.Avatar.AllowedTypes)
	}

	// 确保存储目录存在
	if err := os.MkdirAll(s.storageConfig.Avatar.Path, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %v", err)
	}

	// 生成唯一的文件名
	ext := filepath.Ext(user.Avatar.Filename)
	if ext == "" {
		// 根据文件类型设置默认扩展名
		switch fileType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		}
	}
	fileName := fmt.Sprintf("avatar_%d_%d%s", user.UserID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(s.storageConfig.Avatar.Path, fileName)

	// 保存文件到本地存储
	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		return nil, fmt.Errorf("保存头像文件失败: %v", err)
	}

	// 生成访问URL
	avatarURL := s.storageConfig.Avatar.URL + "/" + fileName

	// 更新用户头像URL
	err = s.userRepository.UpdateUserAvatar(ctx, user.UserID, avatarURL)
	if err != nil {
		// 如果更新失败,删除已上传的文件
		os.Remove(filePath)
		return nil, fmt.Errorf("更新用户头像URL失败: %v", err)
	}

	return &dto.UploadAvatarResponse{
		Code: 200,
		URL:  avatarURL,
	}, nil
}
