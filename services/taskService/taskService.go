package taskService

import (
	"github.com/AtaskTracker/AtaskAPI/database/taskRep"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/AtaskTracker/AtaskAPI/services/googleCloudService"
	"time"
)

type TaskService struct {
	taskRep            *taskRep.TaskRep
	googleCloudService *googleCloudService.GoogleCloudService
}

func New(rep *taskRep.TaskRep) *TaskService {
	return &TaskService{taskRep: rep}
}

const dateFormat = "2006-01-02"

func (s *TaskService) CreateTask(task *dto.Task, userId string) (*dto.Task, error) {
	task.Participants = append(task.Participants, userId)
	url, err2 := s.googleCloudService.UploadImage(task.UUID.Hex(), task.Photo)
	if err2 != nil {
		return task, err2
	}
	task.Photo = url
	var addedTask, err = s.taskRep.CreateTask(*task)
	if err != nil {
		return nil, err
	}
	return &addedTask, nil
}

func (s *TaskService) GetByUserId(userId string) ([]dto.Task, error) {
	var tasks, err = s.taskRep.GetByUserId(userId)
	return tasks, err
}

func (s *TaskService) DeleteById(id string) error {
	return s.taskRep.DeleteById(id)
}

func (s *TaskService) GetTasks(userId string, dateToString string, dateFromString string, label string) ([]dto.Task, error) {
	var dateTo time.Time
	var dateFrom time.Time
	dateTo, err := time.Parse(dateFormat, dateToString)
	if err != nil && dateToString != "" {
		return nil, err
	}
	dateFrom, err = time.Parse(dateFormat, dateFromString)
	if err != nil && dateFromString != "" {
		return nil, err
	}
	return s.taskRep.GetWithFilter(userId, dateTo, dateFrom, label)
}
