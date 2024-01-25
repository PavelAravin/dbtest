package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content,
			author_id,
			assigned_id)
		VALUES ($1, $2,$3,$4) RETURNING id;
		`,
		t.Title,
		t.Content,
		t.AuthorID,
		t.AssignedID,
	).Scan(&id)
	return id, err
}

func (s *Storage) GetAllTasks() ([]Task, error) {

	rows, err := s.db.Query(context.Background(), `		SELECT 			
	id,
	opened,
	closed,
	author_id,
	assigned_id,
	title,
	content 
	FROM tasks;`)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Storage) GetTasksByAuthor(AuthorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `		SELECT	id,	opened,	closed,	author_id,	assigned_id,	title, content	FROM tasks WHERE author_id = $1;`, AuthorID)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil

}

func (s *Storage) GetTasksByAssigned(AssignedID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `		SELECT	id,	opened,	closed,	author_id,	assigned_id,	title, content	FROM tasks WHERE assigned_id = $1;`, AssignedID)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil

}

func (s *Storage) GetTasksByLabel(LabelID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `		SELECT	id,	opened,	closed,	author_id,	assigned_id,	title, content	FROM tasks INNER JOIN tasks_labels ON tasks.id = tasks_labels.task_id  WHERE label_id = $1;`, LabelID)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)

	}
	return tasks, nil
}

func (s *Storage) UpdateTaskByID(ID int, t Task) error {
	_, err := s.db.Exec(context.Background(), `		UPDATE tasks SET	assigned_id = $1,	title = $2,	content = $3 WHERE id = $4;`, t.AssignedID, t.Title, t.Content, ID)
	if err != nil {
		return err
	}
}

func (s *Storage) DeleteTaskByID(ID int) error {
	_, err := s.db.Exec(context.Background(), `		DELETE FROM tasks WHERE id = $1;`, ID)
	if err != nil {
		return err
	}

}

func (s *Storage) GetTaskByID(ID int) (Task, error) {
	var t Task
	row, err := s.db.Query(context.Background(), `		SELECT	id,	opened, closed,	author_id,	assigned_id,	title, content	FROM tasks WHERE id = $1;`, ID)
	if err != nil {
		return t, err
	} else {
		err = row.Scan(&t.ID, &t.Opened, &t.Closed, &t.AuthorID, &t.AssignedID, &t.Title, &t.Content)
		if err != nil {
			return t, err
		}
		return t, nil
	}
}
