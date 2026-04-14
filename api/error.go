package api

import (
	"errors"
	"net/http"

	"github.com/nhxgo/go-utils/errx"
	"github.com/nhxgo/go-utils/logger"
	"github.com/nhxgo/go-utils/request"
)

type ErrorResponse struct {
	Message string `json:"error"`
	Code    string `json:"code,omitempty"`
}

func ErrorJSON(w http.ResponseWriter, ctx *request.Context, err error) {
	if err == nil {
		return
	}

	var apiErr *errx.Error
	if errors.As(err, &apiErr) {
		if ctx != nil {
			ctx.Error(apiErr.Message, "status", apiErr.HttpStatus, "code", string(apiErr.ErrorCode), "error", apiErr.Err)
		} else {
			logger.Error(apiErr.Message, "status", apiErr.HttpStatus, "code", string(apiErr.ErrorCode), "error", apiErr.Err)
		}
		JSON(w, apiErr.HttpStatus, ErrorResponse{Message: apiErr.Message, Code: string(apiErr.ErrorCode)})
	} else {
		JSON(w, http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

}
