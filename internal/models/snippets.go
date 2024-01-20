package models


import (
	"database/sql"
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
	return Snippet{}, nil 
}

func (m *SnippetModel) Latest() ([]Snippet, error){
	latestSnippets := make([]Snippet, 10)
	return latestSnippets, nil 
}
