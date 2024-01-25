package storage

import (
	"dbtest/storage/postgres"
)

type Interface interface {
	Tasks(int, int) ([]postgres.Task, error)
	NewTask(postgres.Task) (int, error)
	UpdateTaskByID(int, postgres.Task) (int, error)
	DeleteTaskByID(int) (int, error)
	GetAllTasks() ([]postgres.Task, error)
	GetTasksByAuthor(int) ([]postgres.Task, error)
	GetTasksByLabel(int) ([]postgres.Task, error)
	GetTasksByAssigned(int) ([]postgres.Task, error)
	GetTaskByID(int) (postgres.Task, error)
}
