package main

import(
	"encoding/json"
	"net/http"

)

func (a *application) healthcheckHandler (w http.ResponseWriter, r *http.Request) {
	data := map[string] string {"status": "available", "environment": a.config.env, "version": a.config.vrs}

	jsResponse, err := json.Marshal(data)
	if err != nil {
    	a.logger.Error(err.Error())
    	http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
    	return
	}

	jsResponse = append(jsResponse, '\n')
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsResponse)

}


// func (app *application) healthcheckHandler (w  http.ResponseWriter, r *http.Request)  {
// 	jsResponse := `{"status": "available", "environment": %q, "version": %q}`
// 	jsResponse = fmt.Sprintf(jsResponse, app.config.env, app.config.vrs)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(jsResponse))
// }