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

// Insert This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// Get This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {

	s := &Snippet{}

	err := m.DB.QueryRow("SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
	//stmt := `SELECT id, title, content, created, expires FROM snippets
	//WHERE expires > UTC_TIMESTAMP() AND id = ?`
	//
	//s := &Snippet{}
	//
	//row := m.DB.QueryRow(stmt, id)
	//
	//// Use row.Scan() to copy the values from each field in sql.Row to the
	//// corresponding field in the Snippet struct. Notice that the arguments
	//// to row.Scan are *pointers* to the place you want to copy the data into,
	//// and the number of arguments must be exactly the same as the number of
	//// columns returned by your statement.
	//err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	//if err != nil {
	//	// If the query returns no rows, then row.Scan() will return a
	//	// sql.ErrNoRows error. We use the errors.Is() function check for that
	//	// error specifically, and return our own ErrNoRecord error instead
	//	if errors.Is(err, sql.ErrNoRows) {
	//		return nil, ErrNoRecord
	//	} else {
	//		return nil, err
	//	}
	//}
	//
	//// If everything went OK then return the Snippet object.
	//return s, nil
}

// Latest This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
