package userHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/AtaskTracker/AtaskAPI/handlers/utilities"
	"github.com/AtaskTracker/AtaskAPI/services/userService"
	"net/http"
)

const contextKeyId = "id"
const tokenCookie = "accessToken"

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

func (h *UserHandler) AuthorizationMW(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(tokenCookie)
		if err != nil {
			utilities.ErrorJsonRespond(w, http.StatusUnauthorized, err)
			return
		}
		user, err := h.userService.Login(&dto.Bearer{Token: tokenCookie.Value})
		if err != nil {
			utilities.ErrorJsonRespond(w, http.StatusUnauthorized, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyId, user.UUID)))
	})
}
