package userService

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/database/userRepo"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"google.golang.org/api/idtoken"
	"os"
)

type UserService struct {
	userRep *userRepo.UserRepo
}

const googleClientId = "954302622465-iruk7dibdhfpl7udjtstl056kaa1sv3e.apps.googleusercontent.com"

func New(rep *userRepo.UserRepo) *UserService {
	return &UserService{userRep: rep}
}

func (s *UserService) Login(bearer *dto.Bearer) (*dto.User, error) {
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

	return &addedUser, nil
}

func (s *UserService) GetUserByEmail(email string) (*dto.User, error) {
	user, err := s.userRep.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func mapToUserDto(payload *idtoken.Payload) *dto.User {
	user := &dto.User{}
	user.Name = payload.Claims["name"].(string)
	user.Email = payload.Claims["email"].(string)
	user.PictureUrl = payload.Claims["picture"].(string)
	return user
}
