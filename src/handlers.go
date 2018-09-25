package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type publicResponse struct {
	Status     string                 `json:"status"`
	Token      string                 `json:"token,omitempty"`
	Message    string                 `json:"message,omitempty"`
	ValidToken map[string]interface{} `json:"token_claims,omitempty"`
	IsValid    bool                   `json:"is_valid,omitempty"`
	Config     *modifiableConfig      `json:"current_config,omitempty"`
}

func handleIssue(w http.ResponseWriter, r *http.Request) {

	response := publicResponse{}
	parseBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	requestDetails := requestToken{}
	err = json.Unmarshal(parseBody, &requestDetails)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = validateTokenRequest(requestDetails)
	if err != nil {
		errorResponse(w, err)
		return
	}

	jwtoken, err := generateToken(requestDetails)
	if err != nil {
		errorResponse(w, err)
		return
	}

	response.Token = jwtoken
	request, _ := json.Marshal(requestDetails)
	successResponse(w, response, "issued token with requested details: "+string(request))
}

func handleValidate(w http.ResponseWriter, r *http.Request) {

	response := publicResponse{}
	parseBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	requestDetails := requestValidate{}
	err = json.Unmarshal(parseBody, &requestDetails)
	if err != nil {
		errorResponse(w, err)
		return
	} else if requestDetails.Token == "" {
		errorResponse(w, errors.New("could not find token"))
		return
	}

	claims, err := validateToken(requestDetails.Token)
	if err != nil {
		errorResponse(w, err)
		return
	}

	response.IsValid = true
	response.ValidToken = claims
	claimsJSON, _ := json.Marshal(claims)
	successResponse(w, response, "validate token with claims: "+string(claimsJSON))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	response := publicResponse{}
	if c.Config {
		response.Status = "running"
	} else {
		response.Status = "not configured"
	}
	if r.Method == http.MethodPost {
		response.Config = c.UserConfig
	}
	c.Logger.Message("success", "server status requested as running")
	output, _ := json.Marshal(response)
	w.Write([]byte(output))
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {

	response := publicResponse{}
	parseBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	requestDetails := modifiableConfig{}
	err = json.Unmarshal(parseBody, &requestDetails)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = updateConfig(requestDetails)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = saveConfig(c.ConfigPath, &c)
	if err != nil {
		errorResponse(w, errors.New("could not persist config"))
		return
	}

	response.Config = c.UserConfig
	successResponse(w, response, "successfully updated config")
}
