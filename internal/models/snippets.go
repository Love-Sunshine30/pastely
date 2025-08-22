package models

import (
	"database/sql"
	"errors"
	"time"
)

// defining a Snippet type to hold individual snippet
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Defining a snippetModel type that wraps around sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// this will insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// sql query for inserting a snippets into the database
	stm := `INSERT INTO snippets(title, content, created, expires)
	VALUES(?, ?, NOW(), DATE_ADD(NOW(), INTERVAL ? DAY))`

	// execute the sql query
	result, err := m.DB.Exec(stm, title, content, expires)
	if err != nil {
		return 0, err
	}

	//get the last inserted snippet's id
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// this will return a specific snippet with specific id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// query for a specific snippet
	stm := `SELECT id, title, content, created, expires
		FROM snippets WHERE expires > NOW() AND id=?`

	// returns a sql.ROW object
	row := m.DB.QueryRow(stm, id)

	//Initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// copy the fields data from sql.Row to s Snippet struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everything is OK then return the Snippet struct
	return s, nil
}

// This will 10 recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// sql query
	stm := `SELECT id, title, content, created, expires
		FROM snippets WHERE expires > NOW() ORDER BY id DESC
		LIMIT 10`

	// Execute the query
	rows, err := m.DB.Query(stm)
	if err != nil {
		return nil, err
	}
	// We use defer to ensure that the above execution is done before the Latest function returns
	defer rows.Close()

	// Initialize a slice to hold 10 latest snippets
	snippets := []*Snippet{}

	for rows.Next() {
		// create a pointer to a new zeroed snippet struct
		s := &Snippet{}

		// Now use rows.Scan() to convert sql rows to snippet struct
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// add the snippet to snippets slice
		snippets = append(snippets, s)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// if everything is OK, return the snippets slice
	return snippets, nil

}
