package main

import (
//   "encoding/json"
  "fmt"
  "net/http"
  // import the data package which contains the definition for Quote

//   "github.com/kelseyaban/qod/internal/data"
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

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}

