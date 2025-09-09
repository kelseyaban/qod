package data

import (
	"time"

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
