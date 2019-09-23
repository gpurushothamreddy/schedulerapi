package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SqliteConfig is the config structure holding information about sqlite db.
type SqliteConfig struct {
	DbName string
}

// Sqlite3Storage is the structure responsible for handling sqlite3 storage.
type SqliteStorage struct {
	config SqliteConfig
	db     *sql.DB
}

// NewSqliteStorage returns a new instance of Sqlite3Storage.
func NewSqliteStorage(config SqliteConfig) SqliteStorage {
	return SqliteStorage{
		config: config,
	}
}

// Connect creates the database file.
func (sqlite *SqliteStorage) Connect() error {
	db, err := sql.Open("sqlite3", sqlite.config.DbName)
	if err != nil {
		return err
	}
	sqlite.db = db
	return nil
}

// Close will close the open DB file.
func (sqlite SqliteStorage) Close() error {
	return sqlite.db.Close()
}

// Initialize will run once to initialize the state of the database into the required one.
// In this case, it'll create the tables required to store task information.
func (sqlite *SqliteStorage) Initialize() error {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS task_store (
        id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        name text,
        params text,
        duration integer,
        last_run text,
        next_run text,
        is_recurring integer,
        hash text
    );
	`
	_, err := sqlite.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}

// Add stores the task to sqlite.
func (sqlite SqliteStorage) Add(task TaskAttributes) error {
	var count int
	rows, err := sqlite.db.Query("SELECT count(*) FROM task_store WHERE hash=?", task.Hash)
	if err == nil {
		rows.Next()
		_ = rows.Scan(&count)
	}
	_ = rows.Close()

	if count == 0 {
		return sqlite.insert(task)
	}
	return nil
}

// Remove will delete the task from storage.
func (sqlite SqliteStorage) Remove(task TaskAttributes) error {
	stmt, err := sqlite.db.Prepare(`DELETE FROM task_store WHERE hash=?`)

	if err != nil {
		return fmt.Errorf("Error while pareparing delete task statement: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		task.Hash,
	)
	if err != nil {
		return fmt.Errorf("Error while deleting task: %s", err)
	}

	return nil
}

// Fetch will return the list of all stored tasks.
func (sqlite SqliteStorage) Fetch() ([]TaskAttributes, error) {
	rows, err := sqlite.db.Query(`
        SELECT name, params, duration, last_run, next_run, is_recurring
        FROM task_store`)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var tasks []TaskAttributes

	for rows.Next() {
		var name, params, lastRun, nextRun, duration, isRecurring string
		err = rows.Scan(&name, &params, &duration, &lastRun, &nextRun, &isRecurring)
		if err != nil {
			return []TaskAttributes{}, err
		}

		task := TaskAttributes{
			Name:        name,
			Params:      params,
			LastRun:     lastRun,
			NextRun:     nextRun,
			Duration:    string(duration),
			IsRecurring: string(isRecurring),
		}

		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return tasks, nil
}

// insert task into task_store table
func (sqlite *SqliteStorage) insert(task TaskAttributes) error {
	stmt, err := sqlite.db.Prepare(`
        INSERT INTO task_store(name, params, duration, last_run, next_run, is_recurring, hash)
        VALUES(?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return fmt.Errorf("Error while pareparing insert task statement: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		task.Name,
		task.Params,
		task.Duration,
		task.LastRun,
		task.NextRun,
		task.IsRecurring,
		task.Hash,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting task: %s", err)
	}

	return nil
}
