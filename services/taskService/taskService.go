package taskService

import (
	"fmt"
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

func (s *TaskService) GetById(taskId string) (*dto.Task, error) {
	task, err := s.taskRep.GetById(taskId)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(task *dto.Task, userId string) error {
	oldTask, err := s.taskRep.GetById(task.UUID.String())
	if err != nil {
		return err
	}
	if oldTask != nil {
		return fmt.Errorf("task not found")
	}
	if !isParticipant(userId, oldTask.Participants) {
		return fmt.Errorf("forbiden: not participant")
	}
	err = s.taskRep.UpdateById(*task)
	return err
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
	percentage := float64(completed) / float64(len(tasks)) * 100
	response := dto.CompletionPercentage{Percentage: percentage}
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
