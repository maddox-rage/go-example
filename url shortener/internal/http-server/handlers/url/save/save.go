package save

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/slg"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"URL" validate:"required, url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

const aliasLenght = 6

func New(log *slog.Logger, URLSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", slg.Err(err))
			render.JSON(w, r, response.Error("failde to decode"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slg.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}
		alias := req.Alias

		if alias == "" {
			alias = random.NewRandomString(aliasLenght)
		}

		id, err := URLSaver.SaveURL(req.URL, alias)

		if errors.Is(err, storage.ErrURLExist) {
			log.Info("url already exist ", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exist"))
			return
		}

		if err != nil {
			log.Error("failed to load URL", slg.Err(err))
			render.JSON(w, r, response.Error("failed to load URL"))
			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}

}
