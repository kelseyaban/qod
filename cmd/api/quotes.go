package main

import (
	//   "encoding/json"
	"fmt"
	"net/http"

	// import the data package which contains the definition for Quote

	"github.com/kelseyaban/qod/internal/data"
	"github.com/kelseyaban/qod/internal/validator"
)

func (a *application)createQuoteHandler(w http.ResponseWriter, r *http.Request) { 
// create a struct to hold a quote
// we use struct tags[``] to make the names display in lowercase
	var incomingData struct {

	Content  string  `json:"content"`
	Author   string  `json:"author"`

	}
    // perform the decoding
	

	// err := json.NewDecoder(r.Body).Decode(&incomingData)
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		// a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	//copy the values from incommingData to a new Quote struct
	//at this point in our code the JSON is well-formed so now we will validate it using the Validator whihc expects a Quote
	quote := &data.Quote {
		Content: incomingData.Content,
		Author: incomingData.Author,
	}

	//INitialize a validator instance
	v := validator.New()

	//do the validation
	data.ValidateQuote(v, quote)
	if !v.IsEmpty() {
		// a.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Add the quote to the database table
	err = a.quoteModel.Insert(quote)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
 
	// Set a Location header. The path to the newly created comment
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/quotes/%d", quote.ID))

	// Send a JSON response with 201 (new resource created) status code
	data := envelope{
		"quote": quote,
	  }
 	err = a.writeJSON(w, http.StatusCreated, data, headers)
 	if err != nil {
	  	a.serverErrorResponse(w, r, err)
	  	return
  	}


 

	
}

