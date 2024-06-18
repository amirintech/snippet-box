package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() AND id = ?`
	snippet := &Snippet{}
	if err := m.DB.QueryRow(stmt, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (m *SnippetModel) GetLatest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []Snippet{}
	for rows.Next() {
		snippet := Snippet{}
		if err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires); err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) Insert(title string, content string, expires int) (*Snippet, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
             VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt = `SELECT id, title, content, created, expires FROM snippets WHERE id = ?`
	snippet := &Snippet{}
	if err := m.DB.QueryRow(stmt, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires); err != nil {
		return nil, err
	}

	return snippet, nil
}
