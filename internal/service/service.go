package service

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"
	cerr "workmate/pkg/customErrors"
	"workmate/pkg/task"

	"github.com/google/uuid"
)

type Config struct {
	AmountOfWorkers int `yaml:"amount_of_workers"`
}

type Service struct {
	tasks           map[uuid.UUID]task.Task
	tasksQueue      []task.Task
	AmountOfWorkers int
	inputTasks      chan task.Task
	outputTasks     chan task.Task
	mx              sync.RWMutex
	isClosed        bool
}

func New(cfg Config) (*Service, error) {
	svc := Service{tasks: make(map[uuid.UUID]task.Task), AmountOfWorkers: cfg.AmountOfWorkers, inputTasks: make(chan task.Task), outputTasks: make(chan task.Task)}
	return &svc, nil
}

func (svc *Service) Start() {
	svc.startPool()
}

func (s *Service) CreateTask() error {
	slog.Info("Service: CreateTask begin")
	id := uuid.New()
	newTask := task.Task{ID: id, Status: task.Created, CreationDate: time.Now()}
	newTask.Context, newTask.CancelFunction = context.WithCancel(context.Background())
	s.mx.Lock()
	s.tasks[id] = newTask
	s.mx.Unlock()
	s.tasksQueue = append(s.tasksQueue, newTask)
	slog.Info("Service: CreateTask: created task", "task", newTask)
	return nil
}

func (s *Service) GetAllTasks() ([]task.Task, error) {
	result := make([]task.Task, 0, len(s.tasks))
	s.mx.RLock()
	for _, t := range s.tasks {
		result = append(result, t)
	}
	s.mx.RUnlock()
	return result, nil
}

func (s *Service) GetTask(id uuid.UUID) (task.Task, error) {
	s.mx.RLock()
	task, ok := s.tasks[id]
	s.mx.RUnlock()
	if !ok {
		return task, cerr.NoSuchTask
	}
	return task, nil
}

func (s *Service) DeleteTask(id uuid.UUID) error {
	s.mx.RLock()
	_, ok := s.tasks[id]
	if !ok {
		return cerr.NoSuchTask
	}
	s.mx.RUnlock()
	s.mx.Lock()
	s.tasks[id].CancelFunction()
	delete(s.tasks, id)
	s.mx.Unlock()
	return nil
}

func (s *Service) stopPool() {
	close(s.inputTasks)
	close(s.outputTasks)
	s.isClosed = true
}

func (s *Service) startPool() {
	s.inputTasks = make(chan task.Task)
	s.outputTasks = make(chan task.Task)

	go func() {
		for range s.AmountOfWorkers {
			go s.startWorker()
		}
	}()

	go func() {
		for {
			if len(s.tasksQueue) > 0 {
				taskToPush := s.tasksQueue[0]
				taskToPush.Status = task.InProgress
				s.tasksQueue = s.tasksQueue[1:len(s.tasksQueue)]
				s.inputTasks <- taskToPush
				s.mx.Lock()
				s.tasks[taskToPush.ID] = taskToPush
				s.mx.Unlock()
			}
		}
	}()

	for finishedTask := range s.outputTasks {
		if finishedTask.Status == task.Cancelled {
			continue
		}
		finishedTask.Status = task.Finished
		s.mx.Lock()
		s.tasks[finishedTask.ID] = finishedTask
		s.mx.Unlock()
	}
}

func (s *Service) startWorker() {
	for newTask := range s.inputTasks {
		start := time.Now().UTC()
		err := simulateWork(newTask.Context)
		if err == cerr.TaskCancelled {
			newTask.Status = task.Cancelled
		}
		newTask.TimeSpent = time.Duration(time.Now().Sub(start).Milliseconds())
		s.outputTasks <- newTask
	}
}

func simulateWork(ctx context.Context) error {
	sleepTime := rand.IntN(3) + 3
	done := make(chan struct{})
	go func() {
		time.Sleep(time.Minute * time.Duration(sleepTime))
		done <- struct{}{}
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return cerr.TaskCancelled
	}
}
