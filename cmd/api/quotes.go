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

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}

