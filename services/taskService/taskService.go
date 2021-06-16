package taskService

import (
	"fmt"
	"github.com/AtaskTracker/AtaskAPI/database/taskRep"
	"github.com/AtaskTracker/AtaskAPI/database/userRepo"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"github.com/AtaskTracker/AtaskAPI/services/googleCloudService"
	"time"
)

type TaskService struct {
	taskRep            *taskRep.TaskRep
	userRep            *userRepo.UserRepo
	googleCloudService *googleCloudService.GoogleCloudService
}

func New(taskRep *taskRep.TaskRep, userRep *userRepo.UserRepo, cloudService *googleCloudService.GoogleCloudService) *TaskService {
	return &TaskService{taskRep: taskRep, userRep: userRep, googleCloudService: cloudService}
}

const dateFormat = "2006-01-02"

func (s *TaskService) CreateTask(task *dto.Task, userId string) (*dto.Task, error) {
	user, err := s.userRep.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	task.Participants = append(task.Participants, user.Email)
	url, err := s.googleCloudService.UploadImage(task.UUID.Hex(), task.Photo)
	if err != nil {
		return task, err
	}
	task.Photo = url
	addedTask, err := s.taskRep.CreateTask(*task)
	if err != nil {
		return nil, err
	}
	return &addedTask, nil
}

func (s *TaskService) GetById(taskId string) (*dto.Task, error) {
	task, err := s.taskRep.GetById(taskId)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(task *dto.Task, userId string) error {
	user, err := s.userRep.GetUserById(userId)
	if err != nil {
		return err
	}
	oldTask, err := s.taskRep.GetById(task.UUID.Hex())
	if err != nil {
		return err
	}
	if oldTask == nil {
		return fmt.Errorf("task not found")
	}
	if !isParticipant(user.Email, oldTask.Participants) {
		return fmt.Errorf("forbiden: not participant")
	}
	err = s.taskRep.UpdateById(*task)
	return err
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
	user, err := s.userRep.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return s.taskRep.GetWithFilter(user.Email, dateTo, dateFrom, label)
}

func (s *TaskService) GetCompletionPercentage(userId string, dateToString string, dateFromString string, label string) (dto.CompletionPercentage, error) {
	tasks, err := s.GetTasks(userId, dateToString, dateFromString, label)
	if err != nil {
		return dto.CompletionPercentage{}, err
	}
	completed := 0
	for _, task := range tasks {
		if task.Status == "done" {
			completed++
		}
	}
	response := dto.CompletionPercentage{Done: completed, Total: len(tasks)}
	return response, nil
}

func (s *TaskService) AddLabel(userId string, taskId string, label dto.Label) error {
	task, err := s.taskRep.GetById(taskId)
	if err != nil {
		return err
	}
	if task != nil {
		return fmt.Errorf("task not found")
	}
	if !isParticipant(userId, task.Participants) {
		return fmt.Errorf("forbiden: not participant")
	}
	return s.taskRep.AddLabel(taskId, label)
}

func isParticipant(userId string, participants []string) bool {
	for _, v := range participants {
		if v == userId {
			return true
		}
	}
	return false
}
