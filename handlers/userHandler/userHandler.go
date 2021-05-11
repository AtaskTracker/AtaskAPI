package userHandler

import (
	"encoding/json"
	"fmt"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/AtaskTracker/AtaskAPI/handlers/utilities"
	"github.com/AtaskTracker/AtaskAPI/services/userService"
	"net/http"
)

type UserHandler struct {
	userService *userService.UserService
}

func New(userService *userService.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Login(writer http.ResponseWriter, request *http.Request) {
	var bearer = &dto.Bearer{}
	if err := json.NewDecoder(request.Body).Decode(bearer); err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusBadRequest, fmt.Errorf("json decode failed"))
		return
	}
	user, err := h.userService.Login(bearer)

	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusCreated, user)
}
