package data

import (
	"time"
	"context"
	"database/sql"

	"errors"
	"github.com/kelseyaban/qod/internal/validator"
)


// each name begins with uppercase so that they are exportable/public
type Quote struct {
	ID        int64     `json:"id"`  // unique value for each quote
	Content   string    `json:"content"`// the quote data
	Author    string    `json:"author"`// the person who wrote the quote
	CreatedAt time.Time `json:"-"`// database timestamp
	Version   int32     `json:"version"`// incremented on each update
}

// create a function that performs the validation check
func ValidateQuote(v *validator.Validator, quote *Quote) {
	//check if the COntent field is empty
	v.Check(quote.Content != "", "content", "cannot be left blank")
	//check if the Author field is empty
	v.Check(quote.Author != "", "author", "cannot be left blank")
	//check if the Content field is empty
	v.Check(len(quote.Content) <= 100, "content", "must no be more than 100 bytes long")
	//check if the AUthor field is empty 
	v.Check(len(quote.Author) <= 25, "author", "must not be more than 25 bytes long")
}

//A QuoteModel expects a connection pool
type QuoteModel struct {
	DB *sql.DB
}

//Insert a new row in the quotes table
//expect a pointer to the actual 
func (q QuoteModel) Insert(quote *Quote) error{

   // the SQL query to be executed against the database table
    query := `INSERT INTO quotes (content, author)
        	VALUES ($1, $2)
        	RETURNING id, created_at, version`

  // the actual values to replace $1, and $2
   args := []any{quote.Content, quote.Author}
  


//Create a context with a 3-second timeout. NO database operation should take more than 3 secs or we will quit it
ctx, cancel := context.WithTimeout(context.Background(), 3  * time.Second)
defer cancel()

// execute the query against the quote database table. We ask for the the id, created_at, and version to be sent back to us which we will use
// to update the Quote struct later on 
return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.ID, &quote.CreatedAt, &quote.Version)

}

// Get a specific Quote from the quotes table
func (q QuoteModel) Get(id int64) (*Quote, error) {
	// check if the id is valid
	 if id < 1 {
		 return nil, ErrRecordNotFound
	 }
	// the SQL query to be executed against the database table
	 query := `
		 SELECT id, created_at, content, author, version
		 FROM quotes
		 WHERE id = $1 `

	 // declare a variable of type Quote to store the returned quote
	 var quote Quote

	 // Set a 3-second context/timer
	 ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	 defer cancel()
	 
	 err := q.DB.QueryRowContext(ctx, query, id).Scan (&quote.ID,&quote.CreatedAt, &quote.Content,&quote.Author, &quote.Version)
	
	// check for which type of error
	if err != nil {
    	switch {
        	case errors.Is(err, sql.ErrNoRows):
            	return nil, ErrRecordNotFound
        	default:
            	return nil, err
        	}
    	}
	return &quote, nil
}

// Update a specific Quotes from the quotes table
func (q QuoteModel) Update(quote *Quote) error {
	// The SQL query to be executed against the database table
	// Every time we make an update, we increment the version number
		query := `
			UPDATE quotes
			SET content = $1, author = $2, version = version + 1
			WHERE id = $3
			RETURNING version`

		args := []any{quote.Content, quote.Author, quote.ID}
   		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
   		defer cancel()

   		return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.Version)

	

}			