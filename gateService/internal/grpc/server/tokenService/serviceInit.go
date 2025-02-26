package tokenService

import (
	"context"
	"gateService/internal/infrastructure/middleware/auth"
)

type Server struct {
	UnimplementedTokenServer
	JWTManager *auth.JWTManager
}

func NewServer(jwtManager *auth.JWTManager) *Server {
	return &Server{
		JWTManager: jwtManager,
	}
}

func (s *Server) TokenVerification(ctx context.Context, in *TokenRequest) (*TokenResponse, error) {
	Claims, err := s.JWTManager.ParseToken(in.Token)
	if err != nil {
		return &TokenResponse{Error: err.Error()}, err
	}

	return &TokenResponse{
		UserID:    int32(Claims.UserInfo.UserID),
		UserName:  Claims.UserInfo.Username,
		Email:     Claims.UserInfo.Email,
		Gender:    Claims.UserInfo.Gender,
		AvatarUrl: Claims.UserInfo.AvatarURL,
		Error:     "",
	}, nil
}
