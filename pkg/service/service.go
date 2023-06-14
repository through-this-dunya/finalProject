package service

import (
	"context"
	"net/http"

	"github.com/through-this-dunya/finalProject.git/pkg/database"
	"github.com/through-this-dunya/finalProject.git/pkg/model"
	"github.com/through-this-dunya/finalProject.git/pkg/proto"
	"github.com/through-this-dunya/finalProject/pkg/database"
	"github.com/through-this-dunya/finalProject/pkg/model"
	"github.com/through-this-dunya/finalProject/pkg/proto"
	"github.com/through-this-dunya/finalProject/pkg/utility"
)

type Server struct {
	Handler database.Handler
	Jwt utility.JwtWrapper
}

func (s *Server) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	var user model.User

	if result := s.Handler.DB.Where(&model.User{Email: req.Email}).First(&user); result.Error == nil {
		return &proto.RegisterResponse{
			Status: http.StatusConflict,
			Error: "Email already exists",
		}, nil
	}

	user.Email = req.Email
	user.Password = utility.HashPassword(req.Password)

	s.Handler.DB.Create(&user)

	return &proto.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user model.User

	if result := s.Handler.DB.Where(&model.User{Email: req.Email}).First(&user); result.Error == nil {
		return &proto.LoginResponse{
			Status: http.StatusNotFound,
			Error: "User not found",
		}, nil
	}

	if !match {
		return &proto.LoginResponse{
			Status: http.StatusNotFound,
			Error: "User not found",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &proto.LoginResponse{
		Status: http.StatusOK,
		Token: token,
	}, nil
}

func (s *Server) Authenticate(ctx context.Context, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	claims, err := s.Jwt.AuthenticateToken(req.Token)

	if err != nil {
		return &proto.AuthenticateResponse{
			Status: http.StatusBadRequest,
			Error: err.Error(),
		}, nil
	}

	var user model.User

	if result := s.Handler.DB.Where(&model.User{Email: claims.Email}).First(&user); result.Error != nil {
		return &proto.AuthenticateResponse{
			Status: http.StatusNotFound,
			Error: "User not found",
		}, nil
	}

	return &proto.AuthenticateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil

}
