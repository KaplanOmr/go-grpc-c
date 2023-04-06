package main

import (
	"context"
	"fmt"
	"go-grpc-c/cache"
	"go-grpc-c/db"
	pb "go-grpc-c/proto"
	"go-grpc-c/utils"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (p *proto) GameResult(ctx context.Context, in *pb.GameResultRequest) (*pb.GameResultResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokenString := strings.Split(md["authorization"][0], " ")[1]

	username, err := utils.CheckJWT(tokenString)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = db.Find("users", map[string]string{"username": username})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if !cache.Set(username, fmt.Sprintf("%f", in.GetScore())) {
		return nil, status.Error(codes.Internal, "SCORE_CANNOT_RECORDED")
	}

	return &pb.GameResultResponse{
		Status:    "success",
		Timestamp: time.Now().UTC().String(),
	}, nil
}

func (p *proto) Leaderboard(ctx context.Context, in *pb.LeaderboardRequest) (*pb.LeaderboardResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokenString := strings.Split(md["authorization"][0], " ")[1]

	username, err := utils.CheckJWT(tokenString)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = db.Find("users", map[string]string{"username": username})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	leaderboard := cache.GetLeaderboard()

	var user User
	var userObject []*pb.UserObject
	for _, userScore := range leaderboard {
		result, err := db.Find("users", map[string]string{"username": userScore.Username})
		if err != nil {
			continue
		}
		result.Decode(&user)

		userObject = append(userObject, &pb.UserObject{
			Username: user.Username,
			Id:       user.ID.String(),
			Score:    userScore.Score,
		})
	}

	return &pb.LeaderboardResponse{
		Status:    "success",
		Timestamp: time.Now().UTC().String(),
		Result:    userObject,
	}, nil
}
