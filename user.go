package main

import (
	"context"
	"go-grpc-c/db"
	pb "go-grpc-c/proto"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	Username string
	Password string
}

func (p *proto) UserRegister(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	if in.GetPassword() == "" || in.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "USERNAME_AND_PASSWORD_REQUIRED")
	}

	user := User{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
	}

	_, err := db.Insert("users", user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserRegisterResponse{
		Status:    "true",
		Timestamp: time.Now().UTC().String(),
		Username:  in.GetUsername(),
		Password:  in.GetPassword(),
	}, nil
}

func (p *proto) UserLogin(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	if in.GetPassword() == "" || in.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "USERNAME_AND_PASSWORD_REQUIRED")
	}

	user := User{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
	}

	_, err := db.Find("users", user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	token, err := generateJWT(user.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserLoginResponse{
		Status:    "true",
		Timestamp: time.Now().UTC().String(),
		Username:  in.GetUsername(),
		Token:     token,
	}, nil
}
