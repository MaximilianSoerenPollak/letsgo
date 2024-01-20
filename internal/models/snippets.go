package models


import (
	"database/sql"
	"errors"
	"time"

)
// =====!!! Note !!!=====
// For testing porupes these return all just empty stuff atm, needs to be fixed. 



//Snippet model, same structure as in the Database table
type Snippet struct {
	ID int
	Title string 
	Content string 
	Created time.Time
	Expires time.Time
}

// Define Snippetmodel type which wraps a sql.Db connection pool? -> Q: Why?
type SnippetModel struct {
	DB *sql.DB 
}

// This will insert a new snippet into the database. 
func (m *SnippetModel) Insert(title, content string, expires int) (int, error){
	// The SQL query but with placehodler values to avoid SQL injection.	

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	
	// For DB.Exec we need to provide the executable query first and then the values of all of the placeholders
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err 
	}
	// LastInsertId gets the ID of the latest inserted Record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err 
	}
	// As id is int64 we need to convert it to be able to return it. Better to do it once here as doing it each time 
	//	we consume it
	return int(id), nil

}

// Return a specific snippet via it's ID 

func (m *SnippetModel) Get(id int) (Snippet, error){
	// We do not want to return any expired Snippets that's why we added a filter. 
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// AS we have the ID as a Primary Key it should only return ALWAYS one row.(or none if it doesn exist)
	row := m.DB.QueryRow(stmt, id)

	// Declare an empty new Snippet struct to use.
	var s Snippet
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)	
	if err != nil {
		// We are checking if we get 0 rows back, as then we can send a better Error message.
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord // This is our own custom error that we have created in the 'errors' file
		} else {
			return Snippet{}, err 
		}
	}
	return s, nil 
}

func (m *SnippetModel) Latest() ([]Snippet, error){
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err 
	}		
	defer rows.Close()
	var latestSnippets []Snippet
	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID,&s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		latestSnippets = append(latestSnippets, s)	
		}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
    // error that was encountered during the iteration. It's important to
    // call this - don't assume that a successful iteration was completed
    // over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err 
	}
	// If everything was okay we go and return.
	return latestSnippets, nil 
}
