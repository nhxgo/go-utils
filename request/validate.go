package request

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nhxgo/go-utils/errx"
)

var validate = validator.New()

func Validate(v any) error {
	if err := validate.Struct(v); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return errx.BadRequestError(err, "invalid request body")
		}

		var sb strings.Builder
		for i, e := range errs {
			str := fmt.Sprintf("%s: %s", strings.ToLower(e.Field()), e.Tag())
			if e.Param() != "" {
				str += fmt.Sprintf("=%s", e.Param())
			}
			sb.WriteString(str)
			if i < len(errs)-1 {
				sb.WriteString(", ")
			}
		}

		return errx.BadRequestError(err, sb.String())
	}
	return nil
}
