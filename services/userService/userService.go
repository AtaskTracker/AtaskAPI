package userService

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/database/labelRep"
	"github.com/AtaskTracker/AtaskAPI/database/userRepo"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/go-redis/redis/v8"
	"google.golang.org/api/idtoken"
	"google.golang.org/appengine/log"
	"os"
	"time"
)

type UserService struct {
	userRep  *userRepo.UserRepo
	redis    *redis.Client
	labelRep *labelRep.LabelRep
}

const googleClientId = "954302622465-iruk7dibdhfpl7udjtstl056kaa1sv3e.apps.googleusercontent.com"

func New(userRep *userRepo.UserRepo, labelRep *labelRep.LabelRep, redis *redis.Client) *UserService {
	return &UserService{userRep: userRep, redis: redis, labelRep: labelRep}
}

func (s *UserService) Login(bearer *dto.Bearer) (*dto.User, error) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "secrets/client_secret.json")
	if err != nil {
		return nil, err
	}
	payload, err := idtoken.Validate(context.Background(), bearer.Token, googleClientId)
	if err != nil {
		return nil, err
	}
	user := mapToUserDto(payload)
	var currentUser dto.User
	existingUser, err := s.userRep.GetUserByEmail(user.Email)
	if !existingUser.UUID.IsZero() {
		user.UUID = existingUser.UUID
		currentUser, err = s.userRep.UpdateUser(*user)
	} else {
		currentUser, err = s.userRep.CreateUser(*user)
	}
	if err != nil {
		return nil, err
	}
	if status := s.redis.Set(context.Background(), bearer.Token, currentUser.UUID.Hex(), time.Hour*24); status.Err() != nil {
		return nil, err
	}
	return &currentUser, nil
}

func (s *UserService) GetUserByEmail(email string) (*dto.User, error) {
	user, err := s.userRep.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &user, nil
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

func (s *UserService) GetLabels(userId string) ([]dto.Label, error) {
	user, err := s.userRep.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	labels, err := s.labelRep.GetLabels(user.Email)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func mapToUserDto(payload *idtoken.Payload) *dto.User {
	user := &dto.User{}
	user.Name = payload.Claims["name"].(string)
	user.Email = payload.Claims["email"].(string)
	user.PictureUrl = payload.Claims["picture"].(string)
	return user
}
