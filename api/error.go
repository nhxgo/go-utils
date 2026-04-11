package api

import (
	"net/http"

	"github.com/nhxgo/go-utils/errx"
	"github.com/nhxgo/go-utils/logger"
	"github.com/nhxgo/go-utils/request"
)

type ErrorResponse struct {
	Message string `json:"error"`
	Code    string `json:"code,omitempty"`
}

func ErrorJSON(w http.ResponseWriter, ctx *request.Context, err errx.Error) {
	if err == nil {
		return
	}
	if ctx != nil {
		ctx.Error(err.Message, "status", err.HttpStatus, "code", string(err.ErrorCode), "error", err.Description)
	} else {
		logger.Error(err.Message, "status", err.HttpStatus, "code", string(err.ErrorCode), "error", err.Description)
	}
	JSON(w, err.HttpStatus, ErrorResponse{Message: err.Message, Code: string(err.ErrorCode)})
}

func Error(w http.ResponseWriter, err errx.Error) {
	if err == nil {
		return
	}
	JSON(w, err.HttpStatus, ErrorResponse{Message: err.Message, Code: string(err.ErrorCode)})
}
