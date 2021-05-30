package userHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/AtaskTracker/AtaskAPI/handlers/utilities"
	"github.com/AtaskTracker/AtaskAPI/services/userService"
	"net/http"
	"strings"
)

type UserHandler struct {
	userService *userService.UserService
}

const contextKeyId = "id"

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

func (h *UserHandler) GetUserByEmail(writer http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")
	if email == "" {
		utilities.ErrorJsonRespond(writer, http.StatusBadRequest, fmt.Errorf("no necessery query param present"))
		return
	}
	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	if user.UUID.IsZero() {
		utilities.ErrorJsonRespond(writer, http.StatusNotFound, fmt.Errorf("no users with given email"))
		return
	}
	utilities.RespondJson(writer, http.StatusCreated, user)
}

func (h *UserHandler) AuthorizationMW(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(reqToken) < 1 {
			utilities.ErrorJsonRespond(w, http.StatusUnauthorized, fmt.Errorf("token not found"))
			return
		}

		reqToken = splitToken[1]

		if reqToken == "test" { // TODO: убрать, сделано для локального тестирования
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyId, "60b3d501385f7aa33124138e")))
			return
		}

		userId, found := h.userService.GetUserId(&dto.Bearer{Token: reqToken})
		if !found {
			utilities.ErrorJsonRespond(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyId, userId)))
	})
}

func (h *UserHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	reqToken := request.Header.Get("Authorization")
	if reqToken == "" {
		utilities.ErrorJsonRespond(writer, http.StatusUnauthorized, fmt.Errorf("token not found"))
		return
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	if err := h.userService.DeleteUserSession(&dto.Bearer{Token: reqToken}); err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusOK, nil)
}

func (h *UserHandler) GetLabels(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(contextKeyId).(string)
	labels, err := h.userService.GetLabels(userId)
	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusOK, labels)
}
