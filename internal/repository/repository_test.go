package repository

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lucsky/cuid"
	"github.com/mtbuzato/go-challenge/internal/model"
	"github.com/stretchr/testify/assert"
)

type CUID struct{}

func (c CUID) Match(v driver.Value) bool {
	err := cuid.IsCuid(v.(string))
	return err == nil
}

func beforeAll(t *testing.T) (*assert.Assertions, *sql.DB, sqlmock.Sqlmock) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub database connection: %s", err)
	}

	return assert, db, mock
}

func TestListAll(t *testing.T) {
	tests := map[string]struct {
		expected []model.Task
		query    func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery
	}{
		"empty": {
			expected: []model.Task{},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return mock.ExpectQuery("SELECT (.+) FROM tasks").
					WillReturnRows(mock.NewRows(nil))
			},
		},
		"one": {
			expected: []model.Task{
				{ID: "1", Name: "Task 1", Completed: false},
			},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return mock.ExpectQuery("SELECT (.+) FROM tasks").
					WillReturnRows(
						mock.NewRows([]string{"id", "name", "completed"}).
							AddRow("1", "Task 1", false),
					)
			},
		},
		"many": {
			expected: []model.Task{
				{ID: "1", Name: "Task 1", Completed: false},
				{ID: "2", Name: "Task 2", Completed: false},
				{ID: "3", Name: "Task 3", Completed: false},
			},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return mock.ExpectQuery("SELECT (.+) FROM tasks").
					WillReturnRows(
						mock.NewRows([]string{"id", "name", "completed"}).
							AddRow("1", "Task 1", false).
							AddRow("2", "Task 2", false).
							AddRow("3", "Task 3", false),
					)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert, db, mock := beforeAll(t)
			defer db.Close()

			test.query(mock)

			repo := NewTaskRepository(db)
			tasks, err := repo.ListAll()

			assert.NoError(err)
			assert.Equal(test.expected, tasks)

			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestListByCompletion(t *testing.T) {
	tests := map[string]struct {
		expected  []model.Task
		query     func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery
		completed bool
	}{
		"empty": {
			expected: []model.Task{},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return mock.ExpectQuery("SELECT (.+) FROM tasks WHERE completed").
					WithArgs(false).
					WillReturnRows(mock.NewRows(nil))
			},
			completed: false,
		},
		"one": {
			expected: []model.Task{
				{ID: "1", Name: "Task 1", Completed: false},
			},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return mock.ExpectQuery("SELECT (.+) FROM tasks WHERE completed").
					WithArgs(false).
					WillReturnRows(
						mock.NewRows([]string{"id", "name", "completed"}).
							AddRow("1", "Task 1", false),
					)
			},
			completed: false,
		},
		"many": {
			expected: []model.Task{
				{ID: "1", Name: "Task 1", Completed: true},
				{ID: "2", Name: "Task 2", Completed: true},
				{ID: "3", Name: "Task 3", Completed: true},
			},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return mock.ExpectQuery("SELECT (.+) FROM tasks WHERE completed").
					WithArgs(true).
					WillReturnRows(
						mock.NewRows([]string{"id", "name", "completed"}).
							AddRow("1", "Task 1", true).
							AddRow("2", "Task 2", true).
							AddRow("3", "Task 3", true),
					)
			},
			completed: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert, db, mock := beforeAll(t)
			defer db.Close()

			test.query(mock)

			repo := NewTaskRepository(db)
			tasks, err := repo.ListByCompletion(test.completed)

			assert.NoError(err)
			assert.Equal(test.expected, tasks)

			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestGetByID(t *testing.T) {
	tests := map[string]struct {
		id       string
		expected model.Task
		query    func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery
	}{
		"invalid_id": {
			id:       "",
			expected: model.Task{},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return nil
			},
		},
		"existing": {
			id: "cl09rb83d000009l13y5n5ur8",
			expected: model.Task{
				ID: "cl09rb83d000009l13y5n5ur8", Name: "Task 1", Completed: false,
			},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id").
					WithArgs("cl09rb83d000009l13y5n5ur8").
					WillReturnRows(
						mock.NewRows([]string{"id", "name", "completed"}).
							AddRow("cl09rb83d000009l13y5n5ur8", "Task 1", false),
					)
			},
		},
		"non_existing": {
			id:       "cl09rb83d000009l13y5n5ur8",
			expected: model.Task{},
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
				return mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id").
					WithArgs("cl09rb83d000009l13y5n5ur8").
					WillReturnRows(
						mock.NewRows(nil),
					)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert, db, mock := beforeAll(t)
			defer db.Close()

			test.query(mock)

			repo := NewTaskRepository(db)
			task, err := repo.GetByID(test.id)

			if test.expected.ID == "" {
				assert.Error(err)
				assert.Empty(task)
			} else {
				assert.NoError(err)
				assert.Equal(test.expected, task)
			}

			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		name        string
		shouldError bool
		query       func(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec
	}{
		"empty_name": {
			name:        "",
			shouldError: true,
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
				return nil
			},
		},
		"valid": {
			name:        "Task 1",
			shouldError: false,
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
				return mock.ExpectExec("INSERT INTO tasks").
					WithArgs(CUID{}, "Task 1", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert, db, mock := beforeAll(t)
			defer db.Close()

			test.query(mock)

			repo := NewTaskRepository(db)
			task, err := repo.Create(test.name)

			if test.shouldError {
				assert.Error(err)
				assert.Empty(task)
			} else {
				assert.NoError(err)
				assert.Equal(test.name, task.Name)
				assert.NoError(cuid.IsCuid(task.ID))
				assert.False(task.Completed)
			}

			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}
func TestUpdate(t *testing.T) {
	tests := map[string]struct {
		task        model.Task
		shouldError bool
		query       func(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec
	}{
		"invalid_id": {
			task: model.Task{
				ID: "",
			},
			shouldError: true,
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
				return nil
			},
		},
		"invalid_name": {
			task: model.Task{
				ID:   "cl09rb83d000009l13y5n5ur8",
				Name: "",
			},
			shouldError: true,
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
				return nil
			},
		},
		"valid": {
			task: model.Task{
				ID:        "cl09rb83d000009l13y5n5ur8",
				Name:      "Task 1",
				Completed: true,
			},
			shouldError: false,
			query: func(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
				return mock.ExpectExec("UPDATE tasks").
					WithArgs("Task 1", true, "cl09rb83d000009l13y5n5ur8").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert, db, mock := beforeAll(t)
			defer db.Close()

			test.query(mock)

			repo := NewTaskRepository(db)
			err := repo.Update(test.task)

			if test.shouldError {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}
