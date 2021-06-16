package taskHandler

import (
	"encoding/json"
	"fmt"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/AtaskTracker/AtaskAPI/handlers/utilities"
	"github.com/AtaskTracker/AtaskAPI/services/taskService"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

const contextKeyId = "id"

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
	userId := request.Context().Value(contextKeyId).(string)

	task, err := h.taskService.CreateTask(task, userId)
	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusCreated, task)
}

func (h *TaskHandler) GetTaskById(writer http.ResponseWriter, request *http.Request) {
	id, _ := mux.Vars(request)["id"]
	// TODO: проверять привелегии, эти всратые
	//contextUserID := request.Context().Value(contextKeyId).(string)
	//if userId != contextUserID {
	//	utilities.ErrorJsonRespond(writer, http.StatusForbidden, fmt.Errorf("no acces for this id"))
	//	return
	//}
	var task, err = h.taskService.GetById(id)
	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	if task == nil {
		utilities.ErrorJsonRespond(writer, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}
	utilities.RespondJson(writer, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(writer http.ResponseWriter, request *http.Request) {
	taskId, _ := mux.Vars(request)["taskId"]
	var task = &dto.Task{}
	if err := json.NewDecoder(request.Body).Decode(task); err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusBadRequest, fmt.Errorf("json decode failed"))
		return
	}
	taskUUID, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}
	task.UUID = taskUUID
	userId := request.Context().Value(contextKeyId).(string)
	err = h.taskService.UpdateTask(task, userId)
	if err != nil {
		switch err.Error() {
		case "task not found":
			utilities.ErrorJsonRespond(writer, http.StatusNotFound, fmt.Errorf("task not found"))
			return
		case "forbiden: not participant":
			utilities.ErrorJsonRespond(writer, http.StatusForbidden, fmt.Errorf("forbiden: not participant"))
			return
		}
	}
	utilities.RespondJson(writer, http.StatusOK, task)
}

func (h *TaskHandler) DeleteById(writer http.ResponseWriter, request *http.Request) {
	id, _ := mux.Vars(request)["id"]
	// TODO: проверять привелегии, эти всратые
	//contextUserID := request.Context().Value(contextKeyId).(string)
	//if userId != contextUserID {
	//	utilities.ErrorJsonRespond(writer, http.StatusForbidden, fmt.Errorf("no acces for this id"))
	//	return
	//}
	if err := h.taskService.DeleteById(id); err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusOK, nil)
}

func (h *TaskHandler) GetUserTasks(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(contextKeyId).(string)
	dateFrom := request.FormValue("dateFrom")
	dateTo := request.FormValue("dateTo")
	label := request.FormValue("label")
	tasks, err := h.taskService.GetTasks(userId, dateTo, dateFrom, label)
	if tasks == nil {
		utilities.ErrorJsonRespond(writer, http.StatusNotFound, fmt.Errorf("no result"))
		return
	}
	if err != nil {
		switch err.(type) {
		default:
			utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
			return
		case *time.ParseError:
			utilities.ErrorJsonRespond(writer, http.StatusBadRequest, err)
			return
		}
	}
	utilities.RespondJson(writer, http.StatusOK, tasks)
}

func (h *TaskHandler) GetCompletionPercentage(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(contextKeyId).(string)
	dateFrom := request.FormValue("dateFrom")
	dateTo := request.FormValue("dateTo")
	label := request.FormValue("label")
	response, err := h.taskService.GetCompletionPercentage(userId, dateTo, dateFrom, label)
	if err != nil {
		switch err.(type) {
		default:
			utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
			return
		case *time.ParseError:
			utilities.ErrorJsonRespond(writer, http.StatusBadRequest, err)
			return
		}
	}
	utilities.RespondJson(writer, http.StatusOK, response)
}

func (h *TaskHandler) AddLabel(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(contextKeyId).(string)
	taskId, _ := mux.Vars(request)["taskId"]
	var label = &dto.Label{}
	if err := json.NewDecoder(request.Body).Decode(label); err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusBadRequest, fmt.Errorf("json decode failed"))
		return
	}
	err := h.taskService.AddLabel(userId, taskId, *label)
	if err != nil {
		switch err.Error() {
		case "task not found":
			utilities.ErrorJsonRespond(writer, http.StatusNotFound, err)
			return
		case "forbiden: not participant":
			utilities.ErrorJsonRespond(writer, http.StatusForbidden, err)
			return
		default:
			utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
			return
		}
	}
	utilities.RespondJson(writer, http.StatusOK, nil)
}
