package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/replicatedhq/kotsadm/pkg/logger"
	"github.com/replicatedhq/kotsadm/pkg/redact"
	"github.com/replicatedhq/kotsadm/pkg/session"
)

type UpdateRedactRequest struct {
	RedactSpec    string `json:"redactSpec"`
	RedactSpecURL string `json:"redactSpecUrl"`
}

type UpdateRedactResponse struct {
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
	UpdatedSpec string `json:"updatedSpec"`
}

type GetRedactResponse struct {
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
	UpdatedSpec string `json:"updatedSpec"`
}

func UpdateRedact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type, origin, accept, authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	updateRedactResponse := UpdateRedactResponse{
		Success: false,
	}

	sess, err := session.Parse(r.Header.Get("Authorization"))
	if err != nil {
		logger.Error(err)
		updateRedactResponse.Error = "failed to parse authorization header"
		JSON(w, 401, updateRedactResponse)
		return
	}

	// we don't currently have roles, all valid tokens are valid sessions
	if sess == nil || sess.ID == "" {
		updateRedactResponse.Error = "failed to parse authorization header"
		JSON(w, 401, updateRedactResponse)
		return
	}

	updateRedactRequest := UpdateRedactRequest{}
	if err := json.NewDecoder(r.Body).Decode(&updateRedactRequest); err != nil {
		logger.Error(err)
		updateRedactResponse.Error = "failed to decode request body"
		JSON(w, 400, updateRedactResponse)
		return
	}

	setSpec := ""
	if updateRedactRequest.RedactSpec != "" {
		setSpec = updateRedactRequest.RedactSpec
	} else if updateRedactRequest.RedactSpecURL != "" {
		resp, err := http.Get(updateRedactRequest.RedactSpecURL)
		if err != nil {
			logger.Error(err)
			updateRedactResponse.Error = "failed to get spec from url"
			JSON(w, 500, updateRedactResponse)
			return
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
			updateRedactResponse.Error = "failed to read spec from url"
			JSON(w, 500, updateRedactResponse)
			return
		}
		setSpec = string(respBytes)
	} else {
		updateRedactResponse.Error = "no spec or url provided"
		JSON(w, 400, updateRedactResponse)
		return
	}

	errMessage, err := redact.SetRedactSpec(setSpec)
	if err != nil {
		logger.Error(err)
		updateRedactResponse.Error = errMessage
		JSON(w, 500, updateRedactResponse)
		return
	}

	updateRedactResponse.Success = true
	updateRedactResponse.UpdatedSpec = setSpec
	JSON(w, 200, updateRedactResponse)
}

func GetRedact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type, origin, accept, authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	getRedactResponse := GetRedactResponse{
		Success: false,
	}

	sess, err := session.Parse(r.Header.Get("Authorization"))
	if err != nil {
		logger.Error(err)
		getRedactResponse.Error = "failed to parse authorization header"
		JSON(w, 401, getRedactResponse)
		return
	}

	// we don't currently have roles, all valid tokens are valid sessions
	if sess == nil || sess.ID == "" {
		getRedactResponse.Error = "failed to parse authorization header"
		JSON(w, 401, getRedactResponse)
		return
	}

	data, errMessage, err := redact.GetRedactSpec()
	if err != nil {
		logger.Error(err)
		getRedactResponse.Error = errMessage
		JSON(w, 500, getRedactResponse)
	}

	getRedactResponse.Success = true
	getRedactResponse.UpdatedSpec = data
	JSON(w, 200, getRedactResponse)
}
