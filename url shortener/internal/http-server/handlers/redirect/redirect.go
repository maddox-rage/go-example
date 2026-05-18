package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, URLGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == ""{
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("not found"))
			return
		}

		resURL, err := URLGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLFound){
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, response.Error("not found"))
			return
		}

		// if err != nil{
		// 	log.Error("failed to get url", slg.Err(err))
		// 	render.JSON(w, r, response.Error("inernal error"))
		// 	return
		// }

		log.Info("got url", slog.String("url", resURL))
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
