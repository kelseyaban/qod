package main

import(
	// "fmt"
	"net/http"
)

func (app *application) healthcheckHandler (w  http.ResponseWriter, r *http.Request)  {
	// panic("Apples & Oranges")

	data := envelope {
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version": app.config.vrs,
		},
	}
	
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r,err)
}

}