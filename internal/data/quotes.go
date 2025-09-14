package data

import (
	"time"
	"context"
	"database/sql"

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

// execute the query against the comments database table. We ask for the the id, created_at, and version to be sent back to us which we will use
// to update the Comment struct later on 
return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.ID, &quote.CreatedAt, &quote.Version)

}