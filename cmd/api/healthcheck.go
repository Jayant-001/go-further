package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	// js := `{"status": "available", "environment": %q, "version": %q}`
	// js = fmt.Sprintf(js, app.config.env, version)

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"status":      "available",
			"environment": app.config.env,
			"version":     version,
		},
	}

	// js, err := json.Marshal(data)
	// if err != nil {
	// 	app.logger.Error(err.Error())
	// 	http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	// 	return
	// }
	// js = append(js, '\n')

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)

	err := app.writeJson(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
