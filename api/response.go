package api

import (
	"encoding/json"
	"net/http"

	"github.com/nhxgo/go-utils/errx"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			ErrorJSON(w, nil, errx.NewError(
				http.StatusInternalServerError,
				errx.UnknownError,
				"Internal Server Error",
				err,
			))
			return
		}

		w.WriteHeader(status)
		_, _ = w.Write(jsonData)
		return
	}

	w.WriteHeader(status)
}
