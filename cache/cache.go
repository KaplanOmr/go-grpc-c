package cache

import (
	"context"
	"log"
	"sort"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type UserScore struct {
	Username string
	Score    float32
}

func Set(username string, score string) bool {
	userScore := map[string]string{"username": username, "score": score}

	for key, value := range userScore {
		err := client.HSet(context.TODO(), "leaderboard:433", key, value).Err()
		if err != nil {
			log.Println(err)
			return false
		}
	}

	return true
}

func GetLeaderboard() []UserScore {
	leaderboard := []UserScore{}

	keys, err := client.Keys(context.TODO(), "leaderboard:*").Result()
	if err != nil {
		return leaderboard
	}

	for _, key := range keys {
		result := client.HGetAll(context.TODO(), key)

		userScoreRaw, err := result.Result()
		if err != nil {
			continue
		}

		score, err := strconv.ParseFloat(userScoreRaw["score"], 32)
		if err != nil {
			continue
		}

		leaderboard = append(leaderboard, UserScore{
			Username: userScoreRaw["username"],
			Score:    float32(score),
		})
	}

	sort.Slice(leaderboard[:], func(i, j int) bool {
		return leaderboard[i].Score < leaderboard[j].Score
	})

	return leaderboard
}
