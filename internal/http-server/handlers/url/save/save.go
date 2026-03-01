package save

import (
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// TODO: доделать
	}

}