package errs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"

	"github.com/bjaus/gokit/vdate"
	"github.com/bjaus/gokit/web"
)

type Response struct {
	Error  string            `json:"error"`
	Detail map[string]string `json:"detail,omitempty"`
}

func Respond(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		web.Respond(w, http.StatusRequestTimeout, Response{Error: "request canceled or deadline exceeded"})
		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		payload := Response{
			Error: "invalid data provided",
			Detail: lo.Associate(ve, func(fe validator.FieldError) (string, string) {
				return fe.Field(), fe.Translate(vdate.Translator())
			}),
		}
		web.Respond(w, http.StatusUnprocessableEntity, payload)
		return
	}

	e := As(err)
	payload := Response{
		Error: e.Message(),
	}

	slog.Error(e.Message(), "error", e.Error())

	web.Respond(w, e.status(), payload)
}
