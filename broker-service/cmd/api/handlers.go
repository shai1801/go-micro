package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"Password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		//DO something
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("Invalid Payload"))

	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	//create json that will be sent to the Authenticate micro-service

	jsonData, _ := json.MarshalIndent(a, "", "\t ")

	//call the microservice

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewReader(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	//get a status code
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer resp.Body.Close()

	//make sure we get back correct status code.

	if resp.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("Invalid Credentials"))
		return

	} else if resp.StatusCode != http.StatusAccepted {

		app.errorJSON(w, errors.New("error calling auth service"))
		return

	}

	//create a variable here to read Body into it

	var jsonFromService jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "authenticated"
	payload.Data = jsonFromService.Data
	app.writeJSON(w, http.StatusAccepted, payload.Data)

}
