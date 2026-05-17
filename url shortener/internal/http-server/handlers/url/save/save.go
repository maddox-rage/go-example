package save

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/slg"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct{
	URL string `json:"URL" validate:"required, url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct{
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface{
	SaveURL(urlToSave string, alias string)(int64, error)
}

func New(log *slog.Logger, URLSaver URLSaver) (http.HandlerFunc){
	 return func(w http.ResponseWriter, r *http.Request){
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err:= render.DecodeJSON(r.Body, &req)

		if err != nil{
			log.Error("failed to decode request body", slg.Err(err))
			render.JSON(w, r, response.Error("failde to decode"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
	 }
}