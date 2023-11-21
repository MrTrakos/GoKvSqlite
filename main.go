package store

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type KVStore struct {
	db *sql.DB
	mu sync.RWMutex
}

func Client(dbPath string) (*KVStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS kv (
            key TEXT PRIMARY KEY,
            value TEXT
        )
    `)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`PRAGMA journal_mode=WAL`)
	if err != nil {
		return nil, err
	}
	return &KVStore{db: db}, nil
}

func (s *KVStore) Set(key string, value interface{}) error {
	_, err := s.db.Exec("INSERT OR REPLACE INTO kv (key, value) VALUES (?, ?)", key, fmt.Sprint(value))
	return err
}

func (s *KVStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var value string
	err := s.db.QueryRow("SELECT value FROM kv WHERE key = ?", key).Scan(&value)
	return value, err
}

func (s *KVStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM kv WHERE key = ?", key)
	return err
}

func (s *KVStore) Exists(key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM kv WHERE key = ?", key).Scan(&count)
	return count > 0, err
}

func (s *KVStore) Keys(pattern string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	var rows *sql.Rows
	var err error

	if pattern == "" || pattern == "*" {
		// If no pattern or "*" is provided, retrieve all keys
		rows, err = s.db.Query("SELECT key FROM kv")
	} else {
		// If a pattern is provided, use it in the SQL query
		rows, err = s.db.Query("SELECT key FROM kv WHERE key LIKE ?", pattern)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

