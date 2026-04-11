package request

import (
	"encoding/json"
	"net/http"

	"github.com/nhxgo/go-utils/errx"
)

func Body(r *http.Request, v any) errx.Error {
	if r.Body == nil {
		return errx.BadRequestError(nil, "request body is required")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return errx.BadRequestError(err, "invalid request body")
	}
	return nil
}
