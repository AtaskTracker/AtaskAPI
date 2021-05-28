package taskService

import (
	"github.com/AtaskTracker/AtaskAPI/database/taskRep"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/AtaskTracker/AtaskAPI/services/googleCloudService"
)

type TaskService struct {
	taskRep            *taskRep.TaskRep
	googleCloudService *googleCloudService.GoogleCloudService
}

func New(rep *taskRep.TaskRep) *TaskService {
	return &TaskService{taskRep: rep}
}

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
