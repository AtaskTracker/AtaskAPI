package userService

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/database/userRepo"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/go-redis/redis/v8"
	"google.golang.org/api/idtoken"
	"google.golang.org/appengine/log"
	"os"
	"time"
)

type UserService struct {
	userRep *userRepo.UserRepo
	redis   *redis.Client
}

func New(rep *userRepo.UserRepo, redis *redis.Client) *UserService {
	return &UserService{userRep: rep, redis: redis}
}

func (s *UserService) Login(bearer *dto.Bearer) (*dto.User, error) {

	const googleClientId = "954302622465-iruk7dibdhfpl7udjtstl056kaa1sv3e.apps.googleusercontent.com"
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "client_secret.json")
	if err != nil {
		return nil, err
	}
	payload, err := idtoken.Validate(context.Background(), bearer.Token, googleClientId)
	if err != nil {
		return nil, err
	}
	user := mapToUserDto(payload)
	addedUser, err := s.userRep.CreateUser(*user)
	if err != nil {
		return nil, err
	}
	if status := s.redis.Set(context.Background(), bearer.Token, addedUser.UUID.String(), time.Hour*24); status.Err() != nil {
		return nil, err
	}

	return &addedUser, nil
}

func (s *UserService) GetUserId(bearer *dto.Bearer) (string, bool) {
	status := s.redis.Get(context.Background(), bearer.Token)
	if status.Err() != nil {
		if status.Err() != redis.Nil {
			log.Warningf(context.Background(), "failed to get id from redis: ", status.Err())
		}
		return "", false
	}
	result, err := status.Result()
	if err != nil {
		log.Warningf(context.Background(), "failed to get id from redis: ", err)
		return "", false
	}
	return result, true
}

func (s *UserService) DeleteUserSession(bearer *dto.Bearer) error {
	status := s.redis.Del(context.Background(), bearer.Token)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func mapToUserDto(payload *idtoken.Payload) *dto.User {
	user := &dto.User{}
	user.Name = payload.Claims["name"].(string)
	user.Email = payload.Claims["email"].(string)
	user.PictureUrl = payload.Claims["picture"].(string)
	return user
}
