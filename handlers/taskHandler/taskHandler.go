package taskHandler

import (
	"encoding/json"
	"fmt"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/AtaskTracker/AtaskAPI/handlers/utilities"
	"github.com/AtaskTracker/AtaskAPI/services/taskService"
	"github.com/gorilla/mux"
	"net/http"
)

const contextKeyUserID = "id"

type TaskHandler struct {
	taskService *taskService.TaskService
}

func New(taskService *taskService.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(writer http.ResponseWriter, request *http.Request) {
	var task = &dto.Task{}
	if err := json.NewDecoder(request.Body).Decode(task); err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusBadRequest, fmt.Errorf("json decode failed"))
		return
	}
	userId := request.Context().Value(contextKeyUserID).(string)

	task, err := h.taskService.CreateTask(task, userId)
	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusCreated, task)
}

func (h *TaskHandler) GetTasksByUserId(writer http.ResponseWriter, request *http.Request) {
	userId, _ := mux.Vars(request)["id"]
	var tasks, err = h.taskService.GetByUserId(userId)
	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusOK, tasks)
}

func (h *TaskHandler) DeleteByUserId(writer http.ResponseWriter, request *http.Request) {
	userId, _ := mux.Vars(request)["id"]
	contextUserID := request.Context().Value(contextKeyUserID).(string)
	if userId != contextUserID {
		utilities.ErrorJsonRespond(writer, http.StatusForbidden, fmt.Errorf("no acces for this id"))
		return
	}
	if err := h.taskService.DeleteById(userId); err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusOK, nil)
}
