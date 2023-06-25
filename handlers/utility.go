package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/polo44/goutilities"
	"polo.gamesmania.io/custodia/gp"
)

// Version is the version of the API
// It is handled by git hooks; should be
const Version string = "0.0.2"

var VersionType *string

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	// ─── OK ─────────────────────────────────────────────────────────────────────────
	w.Header().Set("Content-Type", "application/json")
	goutilities.APIBodyString(&w, `{"result": "All services running"}`)
}

func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	// ─── OK ─────────────────────────────────────────────────────────────────────────
	w.Header().Set("Content-Type", "application/json")
	VersionType = &gp.PConfig.Environment
	goutilities.APIBodyString(&w, fmt.Sprintf(`{"result": "%s-%s"}`, *VersionType, Version))
}

func ReturnJson(w http.ResponseWriter, v interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	bufData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	goutilities.APIBody(&w, &bufData)
	return nil
}

func ReturnJsonOK(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	goutilities.APIBodyString(&w, `{"result": "ok"}`)
}

func ReturnJsonWithCode(w http.ResponseWriter, code int, v interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	bufData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	goutilities.APIBody(&w, &bufData)
	return nil
}
