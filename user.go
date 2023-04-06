package main

import (
	"context"
	"go-grpc-c/db"
	pb "go-grpc-c/proto"
	"go-grpc-c/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}

func (p *proto) UserRegister(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	if in.GetPassword() == "" || in.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "USERNAME_AND_PASSWORD_REQUIRED")
	}

	user := User{
		ID:       primitive.NewObjectID(),
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

	var user User
	result, err := db.Find("users", map[string]string{"username": in.GetUsername(), "password": in.GetPassword()})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result.Decode(&user)

	token, err := utils.GenerateJWT(user.Username)
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
