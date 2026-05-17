package redirect

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, URLURLGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
	}
}
