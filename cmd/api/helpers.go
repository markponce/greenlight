package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelop map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelop, headers http.Header) error {
	// convert to json
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// write new line for command line viewing
	js = append(js, '\n')

	// // Add headers to reponse writer header
	// maps.Insert(w.Header(), maps.All(headers))

	// Add headers to reponse writer header
	for key, value := range headers {
		w.Header()[key] = value
	}

	// add standard header
	w.Header().Set("Content-Type", "application/json")
	// add status code
	w.WriteHeader(status)
	// write json to writer
	w.Write(js)

	// no error
	return nil
}
